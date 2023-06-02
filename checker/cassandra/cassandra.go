package cassandra

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
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
	return "cassandra", nil
}

func (c *Cassandra) Check(ctx context.Context) error {
	cluster := gocql.NewCluster(c.dnss...)
	cluster.Keyspace = "wait4x"

	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		return err
	}

	iter := session.Query("select cql_version from system.local").
		WithContext(ctx).
		Iter()
	defer iter.Close()

	rows, err := iter.RowData()
	if err != nil {
		return err
	}

	if len(rows.Values) != 1 {
		return ErrBadConnection
	}

	if ok := iter.Scan(rows.Values...); !ok {
		return ErrBadConnection
	}

	values, ok := rows.Values[0].(*string)
	if !ok {
		return ErrBadConnection
	}

	if len(*values) < 1 {
		return ErrBadConnection
	}

	return nil
}
