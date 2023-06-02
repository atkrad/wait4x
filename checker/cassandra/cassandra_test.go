package cassandra

import (
	"context"
	"sync"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

var containerOnce sync.Once

func getCassandraContiner(ctx context.Context, t *testing.T) testcontainers.Container {
	containerOnce.Do(func() {
		req := testcontainers.ContainerRequest{
			Image: "cassandra:latest",
			ExposedPorts: []string{
				"9042:9042",
			},
			Env: map[string]string{
				"CASSANDRA_CLUSTER_NAME": "wait4x_cluster",
			},
		}
	})
}
