package cassandra

import (
	"context"
	"errors"

	"wait4x.dev/v2/checker"
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
	return "", errors.New("not impl yet")
}

func (c *Cassandra) Check(ctx context.Context) error {
	return errors.New("not impl yet")
}
