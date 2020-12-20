package checker

import "github.com/atkrad/wait4x/pkg/log"

type Checker interface {
	SetLogger(logger log.Logger)
	Check() bool
}
