// Copyright 2020 Mohammad Abdolirad
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checker

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

// Redis represents Redis checker
type Redis struct {
	address   string
	expectKey string
	timeout   time.Duration
	*LogAware
}

// NewRedis creates the Redis checker
func NewRedis(address string, expectKey string, timeout time.Duration) Checker {
	r := &Redis{
		address:   address,
		expectKey: expectKey,
		timeout:   timeout,
		LogAware:  &LogAware{},
	}

	return r
}

// Check checks Redis connection
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
