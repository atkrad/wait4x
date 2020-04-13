package cmd

import (
	"errors"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
	"time"
)

var redisCmd = &cobra.Command{
	Use:   "redis ADDRESS",
	Short: "Check Redis connection or key existence.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("ADDRESS is required argument for the redis command")
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
	Run: func(cmd *cobra.Command, args []string) {
		ticker := time.NewTicker(Interval)
		defer ticker.Stop()

		go func() {
			connectionTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
			password, _ := cmd.Flags().GetString("password")
			db, _ := cmd.Flags().GetInt("db")
			expectKey, _ := cmd.Flags().GetString("expect-key")

			for ; true; <-ticker.C {
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

					continue
				} else {
					// It can connect to Redis successfully
					if expectKey == "" {
						os.Exit(EXIT_SUCCESS)
					}

					splittedKey := strings.Split(expectKey, "=")
					keyHasValue := len(splittedKey) == 2

					val, err := client.Get(splittedKey[0]).Result()
					if err == redis.Nil {
						// Redis key does not exist.
						log.WithFields(log.Fields{
							"key": splittedKey[0],
						}).Info("Key does not exist.")

						continue
					} else if err != nil {
						// Error occurred on get Redis key
						log.Debug(err)

						continue
					} else {
						// The Redis key exists.
						if !keyHasValue {
							os.Exit(EXIT_SUCCESS)
						}

						// When the user expect a key with value
						matched, _ := regexp.MatchString(splittedKey[1], val)
						if matched {
							os.Exit(EXIT_SUCCESS)
						} else {
							log.WithFields(log.Fields{
								"key":    splittedKey[0],
								"actual": val,
								"expect": splittedKey[1],
							}).Info("Checking value expectation of the key")

							continue
						}
					}
				}
			}
		}()

		time.Sleep(Timeout)
		log.Info("Operation Timed Out")

		os.Exit(EXIT_TIMEDOUT)
	},
}

func init() {
	rootCmd.AddCommand(redisCmd)

	redisCmd.Flags().Duration("connection-timeout", time.Second*5, "Dial timeout for establishing new connections.")
	redisCmd.Flags().String("password", "", "Optional password. Must match the password specified in the requirepass server configuration option.")
	redisCmd.Flags().String("expect-key", "", "Checking key existence.")
	redisCmd.Flags().Int("db", 0, "Database to be selected after connecting to the server.")
}
