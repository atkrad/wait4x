package main

import (
	"github.com/atkrad/wait4x/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	cmd.Execute()
}
