package checker

import (
	"regexp"
	"strings"
	"time"

	"github.com/atkrad/wait4x/pkg/log"
	"github.com/go-redis/redis/v7"
)

type Redis struct {
	logger    log.Logger
	address   string
	expectKey string
	timeout   time.Duration
}

func NewRedis(address string, expectKey string, timeout time.Duration) Checker {
	r := &Redis{
		address:   address,
		expectKey: expectKey,
		timeout:   timeout,
	}

	return r
}

func (r *Redis) SetLogger(logger log.Logger) {
	r.logger = logger
}

func (r *Redis) Check() bool {
	r.logger.Info("Checking Redis connection ...")

	opts, err := redis.ParseURL(r.address)
	if err != nil {
		r.logger.Debug(err)

		return false
	}
	opts.DialTimeout = r.timeout

	client := redis.NewClient(opts)

	// Check Redis connection
	_, err = client.Ping().Result()
	if err != nil {
		r.logger.Debug(err)

		return false
	}

	// It can connect to Redis successfully
	if r.expectKey == "" {
		return true
	}

	splittedKey := strings.Split(r.expectKey, "=")
	keyHasValue := len(splittedKey) == 2

	val, err := client.Get(splittedKey[0]).Result()
	if err == redis.Nil {
		// Redis key does not exist.
		r.logger.InfoWithFields("Key does not exist.", map[string]interface{}{"key": splittedKey[0]})

		return false
	}

	if err != nil {
		// Error occurred on get Redis key
		r.logger.Debug(err)

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

	r.logger.InfoWithFields("Checking value expectation of the key", map[string]interface{}{"key": splittedKey[0], "actual": val, "expect": splittedKey[1]})

	return false
}
