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

package redis

import (
	"context"
	"github.com/atkrad/wait4x/pkg/checker"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

var ParseURLErr = checker.NewError("parse url error", "debug")
var PingErr = checker.NewError("ping error", "debug")
var KeyExistenceErr = checker.NewError("the key doesn't exist", "info")
var GetKeyErr = checker.NewError("get key", "debug")
var KeyValueExistenceErr = checker.NewError("the key and desired value doesn't exist", "info")

// Option configures a Redis.
type Option func(s *Redis)

// Redis represents Redis checker
type Redis struct {
	address   string
	expectKey string
	timeout   time.Duration
	*checker.LogAware
}

// New creates the Redis checker
func New(address string, opts ...Option) checker.Checker {
	r := &Redis{
		address:  address,
		timeout:  time.Second * 5,
		LogAware: &checker.LogAware{},
	}

	// apply the list of options to Redis
	for _, opt := range opts {
		opt(r)
	}

	return r
}

// WithTimeout configures a timeout for establishing new connections
func WithTimeout(timeout time.Duration) Option {
	return func(r *Redis) {
		r.timeout = timeout
	}
}

// WithExpectKey configures a key expectation
func WithExpectKey(key string) Option {
	return func(r *Redis) {
		r.expectKey = key
	}
}

// Check checks Redis connection
func (r *Redis) Check(ctx context.Context) error {
	r.Logger().Info("Checking Redis connection ...")

	opts, err := redis.ParseURL(r.address)
	if err != nil {
		return ParseURLErr.WithWrap(err)
	}
	opts.DialTimeout = r.timeout

	client := redis.NewClient(opts)

	// Check Redis connection
	_, err = client.WithContext(ctx).Ping().Result()
	if err != nil {
		return PingErr.WithWrap(err)
	}

	// It can connect to Redis successfully
	if r.expectKey == "" {
		return nil
	}

	splittedKey := strings.Split(r.expectKey, "=")
	keyHasValue := len(splittedKey) == 2

	val, err := client.WithContext(ctx).Get(splittedKey[0]).Result()
	if err == redis.Nil {
		// Redis key does not exist.
		return KeyExistenceErr.WithFields(map[string]interface{}{"key": splittedKey[0]})
	}

	if err != nil {
		// Error occurred on get Redis key
		return GetKeyErr.WithWrap(err)
	}

	// The Redis key exists and user doesn't want to match value
	if !keyHasValue {
		return nil
	}

	// When the user expect a key with value
	matched, _ := regexp.MatchString(splittedKey[1], val)
	if matched {
		return nil
	}

	return KeyValueExistenceErr.WithFields(map[string]interface{}{"key": splittedKey[0], "actual": val, "expect": splittedKey[1]})
}
