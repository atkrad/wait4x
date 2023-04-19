package redis

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"sync"
	"testing"
	"time"
)

var containerOnce sync.Once
var redisContainer testcontainers.Container

func getRedisContainer(ctx context.Context, t *testing.T) testcontainers.Container {
	containerOnce.Do(func() {
		req := testcontainers.ContainerRequest{
			Image:        "redis:latest",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		}

		var err error
		redisContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           testcontainers.TestLogger(t),
		})
		if err != nil {
			t.Fatal(err)
		}

		//t.Cleanup(func() {
		//	if err := redisContainer.Terminate(ctx); err != nil {
		//		t.Fatalf("failed to terminate container: %s", err)
		//	}
		//})
	})

	return redisContainer
}

func TestValidAddress(t *testing.T) {
	ctx := context.Background()
	rc := getRedisContainer(ctx, t)

	endpoint, err := rc.Endpoint(ctx, "redis")
	if err != nil {
		t.Fatal(err)
	}

	checker := New(endpoint)
	assert.Nil(t, checker.Check(ctx))
}

func TestKeyExistence(t *testing.T) {
	ctx := context.Background()
	redisContainer := getRedisContainer(ctx, t)

	endpoint, err := redisContainer.Endpoint(ctx, "redis")
	if err != nil {
		t.Fatal(err)
	}

	opts, err := redis.ParseURL(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	redisClient := redis.NewClient(opts)
	redisClient.Set("Foo", "Bar", time.Hour)

	checker := New(endpoint, WithExpectKey("Foo"))
	assert.Nil(t, checker.Check(ctx))
}
