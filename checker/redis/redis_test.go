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
	"wait4x.dev/v2/checker"
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

func TestInvalidConnection(t *testing.T) {
	var expectedError *checker.ExpectedError
	chk := New("redis://127.0.0.1:8787")
	assert.ErrorAs(t, chk.Check(context.Background()), &expectedError)
}

func TestValidAddress(t *testing.T) {
	ctx := context.Background()
	rc := getRedisContainer(ctx, t)

	endpoint, err := rc.Endpoint(ctx, "redis")
	if err != nil {
		t.Fatal(err)
	}

	chk := New(endpoint)
	assert.Nil(t, chk.Check(ctx))
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

	chk := New(endpoint, WithExpectKey("Foo"))
	assert.Nil(t, chk.Check(ctx))

	chk = New(endpoint, WithExpectKey("Foo=^B.*$"))
	assert.Nil(t, chk.Check(ctx))

	var expectedError *checker.ExpectedError
	chk = New(endpoint, WithExpectKey("Foo=^b[A-Z]$"))
	assert.ErrorAs(t, chk.Check(ctx), &expectedError)

	chk = New(endpoint, WithExpectKey("Bob"))
	assert.ErrorAs(t, chk.Check(ctx), &expectedError)
}
