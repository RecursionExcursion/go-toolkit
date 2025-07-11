package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection[T any] struct {
	Db         string
	Collection string
	atlasUrI   string
}

/* Generic Mongo Query Fns */

// Read
func (mc *MongoConnection[T]) FindTById(id string) (T, error) {
	query := bson.D{primitive.E{
		Key:   "id",
		Value: id,
	}}

	return mongoQuery(*mc, func(c *mongo.Collection) (T, error) {
		var t T
		res := c.FindOne(context.TODO(), query)

		err := res.Decode(&t)
		if err == mongo.ErrNoDocuments {
			log.Println("No document found")
			return t, err
		}
		if err != nil {
			log.Println("Error decoding from Mongo:", err)
			return t, err
		}
		return t, nil
	})
}

func (mc *MongoConnection[T]) FindFirstT() (T, error) {
	return mongoQuery(
		*mc, func(c *mongo.Collection) (T, error) {
			res := c.FindOne(context.TODO(), bson.M{})

			var t T
			err := res.Decode(&t)
			if err == mongo.ErrNoDocuments {
				log.Println("No document found", err)
				return t, err
			}
			if err != nil {
				log.Println("Error decoding from Mongo", err)
				return t, err
			}

			return t, nil
		},
	)
}

// Create/Update
/* SaveT- Always inserts the document */
func (mc *MongoConnection[T]) SaveT(t T) (bool, error) {
	_, err := mongoQuery(*mc, func(c *mongo.Collection) (T, error) {
		_, err := c.InsertOne(context.TODO(), t)
		if err != nil {
			return t, err
		}
		return t, nil
	})
	return err == nil, err
}

/* UpsertT- Updates doc if filter matches, otherwise it creates a new entry */
func (mc *MongoConnection[T]) UpsertT(t T, filter bson.M) (bool, error) {
	_, err := mongoQuery(*mc, func(c *mongo.Collection) (T, error) {
		opts := options.Update().SetUpsert(true)

		// Perform upsert: update matching doc or insert if not exists
		_, err := c.UpdateOne(context.TODO(), filter, bson.M{"$set": t}, opts)
		if err != nil {
			return t, err
		}
		return t, nil
	})
	return err == nil, err
}

// Delete
func (mc *MongoConnection[T]) DeleteById(id string) (bool, error) {
	var t T
	_, err := mongoQuery(*mc, func(c *mongo.Collection) (T, error) {
		query := bson.D{primitive.E{
			Key:   "id",
			Value: id,
		}}
		_, err := c.DeleteOne(context.Background(), query)
		if err != nil {
			return t, err
		}
		return t, nil
	})
	return err == nil, err
}

// connetion logic
func connectMongoClient(atlasUri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(atlasUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type queryFn[T any] func(c *mongo.Collection) (T, error)

/* mongoQuery creates, manages, and closes connection while executingthe custom query fn passed in */
func mongoQuery[T any](mc MongoConnection[T], query queryFn[T]) (T, error) {
	client, err := connectMongoClient(mc.atlasUrI)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database(mc.Db)
	c := db.Collection(mc.Collection)
	return query(c)
}
