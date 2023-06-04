package cassandra

import (
	"context"
	"errors"
	"github.com/gocql/gocql"
	"strings"
	"wait4x.dev/v2/checker"
)

var (
	ErrBadConnection = errors.New("bad connection")
)

type ConnectionParams struct {
	Hosts    []string
	Username string
	Password string
}

type Cassandra struct {
	connectionParams ConnectionParams
}

func New(connectionParams ConnectionParams) checker.Checker {
	return &Cassandra{
		connectionParams: connectionParams,
	}
}

func (c *Cassandra) Identity() (string, error) {
	return strings.Join(c.connectionParams.Hosts, ","), nil
}

func (c *Cassandra) Check(ctx context.Context) error {
	session, err := c.connectToCluster()
	if err != nil {
		return checker.NewExpectedError(
			"failed to establish a connection to the cassandra cluster",
			err,
			"connection",
			c.connectionParams.Hosts,
		)
	}
	defer session.Close()

	iter := session.Query("select cql_version from system.local").
		WithContext(ctx).
		Iter()
	defer iter.Close()

	rows, err := iter.RowData()
	if err != nil {
		return checker.NewExpectedError(
			"failed to get the row data",
			err,
			"rowData",
			c.connectionParams.Hosts,
		)
	}

	if len(rows.Values) != 1 {
		return checker.NewExpectedError(
			"failed to query system.local",
			ErrBadConnection,
			"values",
			c.connectionParams.Hosts,
		)
	}

	if ok := iter.Scan(rows.Values...); !ok {
		return checker.NewExpectedError(
			"failed to scan row values",
			ErrBadConnection,
			"scan",
			c.connectionParams.Hosts,
		)
	}

	values, ok := rows.Values[0].(*string)
	if !ok {
		return checker.NewExpectedError(
			"failed to convert scanned values",
			ErrBadConnection,
			"conversion",
			c.connectionParams.Hosts,
		)
	}

	if len(*values) < 1 {
		return checker.NewExpectedError(
			"no returning values",
			ErrBadConnection,
			"return",
			c.connectionParams.Hosts,
		)
	}

	return nil
}

func (c *Cassandra) connectToCluster() (*gocql.Session, error) {
	cluster := gocql.NewCluster(c.connectionParams.Hosts...)
	cluster.Keyspace = "system"
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.All

	if c.connectionParams.Username != "" && c.connectionParams.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: c.connectionParams.Username,
			Password: c.connectionParams.Password,
		}
	}

	return cluster.CreateSession()
}
