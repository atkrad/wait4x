package checker

import "github.com/atkrad/wait4x/pkg/log"

// Checker is the interface that wraps the basic checker methods.
type Checker interface {
	SetLogger(logger log.Logger)
	Check() bool
}
