package cassandra

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"sync"
	"testing"
	"wait4x.dev/v2/checker"
)

var containerOnce sync.Once
var cassandraContainer testcontainers.Container

func getCassandraContainer(ctx context.Context, t *testing.T) testcontainers.Container {
	containerOnce.Do(func() {
		req := testcontainers.ContainerRequest{
			Image: "cassandra:latest",
			ExposedPorts: []string{
				"9042/tcp",
				"7000-7001/tcp",
				"7199/tcp",
				"9160/tcp",
			},

			Name: "cassandra-instance-1",
			Env: map[string]string{
				"CASSANDRA_CLUSTER_NAME":      "wait4x_cluster",
				"CASSANDRA_BROADCAST_ADDRESS": "127.0.0.1",
			},
			WaitingFor: wait.ForLog("Created default superuser role 'cassandra'"),
		}

		var err error
		cassandraContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           testcontainers.TestLogger(t),
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	return cassandraContainer
}

func TestInvalidConnection(t *testing.T) {
	var exceptError *checker.ExpectedError

	chk := New(ConnectionParams{
		Hosts: []string{"127.0.0.1:9042"},
	})

	assert.ErrorAs(t, chk.Check(context.Background()), &exceptError)
}

func TestValidConnectionSingleNode(t *testing.T) {
	ctx := context.Background()
	rc := getCassandraContainer(ctx, t)

	mappedPort, err := rc.MappedPort(ctx, "9042/tcp")
	if err != nil {
		t.Fatal(err)
	}

	cassURL := fmt.Sprintf("127.0.0.1:%d", mappedPort.Int())
	chk := New(ConnectionParams{
		Hosts: []string{cassURL},
	})
	chkErr := chk.Check(ctx)

	assert.Nil(t, chkErr)
}
