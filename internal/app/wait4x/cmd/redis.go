package cmd

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewRedisCommand creates the redis sub-command
func NewRedisCommand() *cobra.Command {
	redisCommand := &cobra.Command{
		Use:   "redis ADDRESS",
		Short: "Check Redis connection or key existence.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("ADDRESS is required argument for the redis command")
			}

			return nil
		},
		Example: `
  # Checking Redis connection
  wait4x redis 127.0.0.1:6379

  # Checking a key existence
  wait4x redis 127.0.0.1:6379 --expect-key FOO

  # Checking a key existence and matching the value
  wait4x redis 127.0.0.1:6379 --expect-key "FOO=^b[A-Z]r$"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			for !checkingRedis(cmd, args) {
				select {
				case <-ctx.Done():
					return errors.NewTimedOutError()
				case <-time.After(Interval):
				}
			}

			return nil
		},
	}

	redisCommand.Flags().Duration("connection-timeout", time.Second*5, "Dial timeout for establishing new connections.")
	redisCommand.Flags().String("password", "", "Optional password. Must match the password specified in the requirepass server configuration option.")
	redisCommand.Flags().String("expect-key", "", "Checking key existence.")
	redisCommand.Flags().Int("db", 0, "Database to be selected after connecting to the server.")

	return redisCommand
}

func checkingRedis(cmd *cobra.Command, args []string) bool {
	connectionTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
	password, _ := cmd.Flags().GetString("password")
	db, _ := cmd.Flags().GetInt("db")
	expectKey, _ := cmd.Flags().GetString("expect-key")

	log.Info("Checking Redis connection ...")

	client := redis.NewClient(&redis.Options{
		Addr:        args[0],
		Password:    password,
		DB:          db,
		DialTimeout: connectionTimeout,
	})

	// Check Redis connection
	_, err := client.Ping().Result()
	if err != nil {
		log.Debug(err)

		return false
	}

	// It can connect to Redis successfully
	if expectKey == "" {
		return true
	}

	splittedKey := strings.Split(expectKey, "=")
	keyHasValue := len(splittedKey) == 2

	val, err := client.Get(splittedKey[0]).Result()
	if err == redis.Nil {
		// Redis key does not exist.
		log.WithFields(log.Fields{
			"key": splittedKey[0],
		}).Info("Key does not exist.")

		return false
	}

	if err != nil {
		// Error occurred on get Redis key
		log.Debug(err)

		return false
	}

	// The Redis key exists and user doesn't want to match value
	if !keyHasValue {
		return true
	}

	// When the user expect a key with value
	matched, _ := regexp.MatchString(splittedKey[1], val)
	if matched {
		return true
	}

	log.WithFields(log.Fields{
		"key":    splittedKey[0],
		"actual": val,
		"expect": splittedKey[1],
	}).Info("Checking value expectation of the key")

	return false
}
