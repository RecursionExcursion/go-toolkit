package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type EnvLoader struct {
	loaded bool
}

// returns value associated with key or an error if no key is found
func (el *EnvLoader) Get(key string) (string, error) {
	if !el.loaded {
		if err := el.LoadEnv(); err != nil {
			return "", err
		}
	}

	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("env key %v not set", key)
	}
	return val, nil
}

// returns value associated with key or panics if no key is found
func (el *EnvLoader) MustGet(key string) string {
	if !el.loaded {
		if err := el.LoadEnv(); err != nil {
			panic(err)
		}
	}

	val := os.Getenv(key)
	if val == "" {
		log.Panicf("Env key %v not set", key)
	}
	return val
}

// returns value associated with key or fallback if no key is found
func (el *EnvLoader) GetOrFallback(key string, fallback string) string {
	if !el.loaded {
		if err := el.LoadEnv(); err != nil {
			return fallback
		}
	}

	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// returns value associated with key or zero value ("") if no key is found
func (el *EnvLoader) GetOrDefault(key string) string {
	if !el.loaded {
		if err := el.LoadEnv(); err != nil {
			return ""
		}
	}
	return os.Getenv(key)
}

// LoadEnv loads environment variables from one or more files into the process's environment.
// Files should contain lines in the format KEY=VALUE.
// If no files are provided, it defaults to loading from ".env".
func (el *EnvLoader) LoadEnv(files ...string) error {

	if len(files) == 0 {
		return loadOsEnv(".env")
	}

	for _, f := range files {
		if err := loadOsEnv(f); err != nil {
			return err
		}
	}
	el.loaded = true
	return nil
}

func loadOsEnv(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		/* TODO Need to impl logic to cut out trailing comments
		 * '#' should only be permitted inside quotes
		 * if not then everything after needs to be treated as a comment
		 */
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		val = strings.Trim(val, `"'`)

		os.Setenv(key, val)
	}

	return scanner.Err()
}
