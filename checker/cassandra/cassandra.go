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

type Cassandra struct {
	dnss []string
}

func New(dnss []string) checker.Checker {
	return &Cassandra{
		dnss: dnss,
	}
}

func (c *Cassandra) Identity() (string, error) {
	return strings.Join(c.dnss, ","), nil
}

func (c *Cassandra) Check(ctx context.Context) error {
	cluster := gocql.NewCluster(c.dnss...)
	cluster.Keyspace = "system"
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.All

	session, err := cluster.CreateSession()
	if err != nil {
		return checker.NewExpectedError(
			"failed to establish a connection to the cassandra cluster",
			err,
			"connection",
			c.dnss,
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
			c.dnss,
		)
	}

	if len(rows.Values) != 1 {
		return checker.NewExpectedError(
			"failed to query system.local",
			ErrBadConnection,
			"values",
			c.dnss,
		)
	}

	if ok := iter.Scan(rows.Values...); !ok {
		return checker.NewExpectedError(
			"failed to scan row values",
			ErrBadConnection,
			"scan",
			c.dnss,
		)
	}

	values, ok := rows.Values[0].(*string)
	if !ok {
		return checker.NewExpectedError(
			"failed to convert scanned values",
			ErrBadConnection,
			"conversion",
			c.dnss,
		)
	}

	if len(*values) < 1 {
		return checker.NewExpectedError(
			"no returning values",
			ErrBadConnection,
			"return",
			c.dnss,
		)
	}

	return nil
}
