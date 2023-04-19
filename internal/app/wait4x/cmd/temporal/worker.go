// Copyright 2023 The Wait4X Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package temporal

import (
	"errors"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"wait4x.dev/v2/checker/temporal"
	"wait4x.dev/v2/waiter"
)

// NewWorkerCommand creates the worker sub-command
func NewWorkerCommand() *cobra.Command {
	workerCommand := &cobra.Command{
		Use:   "worker TARGET [flags] [-- command [args...]]",
		Short: "Check Temporal worker registration",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("TARGET is required argument for the worker command")
			}

			return nil
		},
		Example: `
  # Checking a task queue that has registered workers (pollers) or not
  wait4x temporal worker 127.0.0.1:7233 --namespace __YOUR_NAMESPACE__ --task-queue __YOUR_TASK_QUEUE__

  # Checking the specific a Temporal worker (pollers)
  wait4x temporal worker 127.0.0.1:7233 --namespace __YOUR_NAMESPACE__ --task-queue __YOUR_TASK_QUEUE__ --expect-worker-identity-regex ".*@__HOSTNAME__@.*"
`,
		RunE: runWorker,
	}

	workerCommand.Flags().String("namespace", "", "Temporal namespace.")
	workerCommand.Flags().String("task-queue", "", "Temporal task queue.")
	workerCommand.Flags().String("expect-worker-identity-regex", "", "Expect Temporal worker (poller) identity regex.")
	cobra.MarkFlagRequired(workerCommand.Flags(), "namespace")
	cobra.MarkFlagRequired(workerCommand.Flags(), "task-queue")

	return workerCommand
}

func runWorker(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	conTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
	insecureTransport, _ := cmd.Flags().GetBool("insecure-transport")
	insecureSkipTLSVerify, _ := cmd.Flags().GetBool("insecure-skip-tls-verify")

	namespace, _ := cmd.Flags().GetString("namespace")
	taskQueue, _ := cmd.Flags().GetString("task-queue")
	expectWorkerIdentityRegex, _ := cmd.Flags().GetString("expect-worker-identity-regex")

	logger, err := logr.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	// ArgsLenAtDash returns -1 when -- was not specified
	if i := cmd.ArgsLenAtDash(); i != -1 {
		args = args[:i]
	}

	tc := temporal.New(
		temporal.CheckModeWorker,
		args[0],
		temporal.WithTimeout(conTimeout),
		temporal.WithInsecureTransport(insecureTransport),
		temporal.WithInsecureSkipTLSVerify(insecureSkipTLSVerify),
		temporal.WithNamespace(namespace),
		temporal.WithTaskQueue(taskQueue),
		temporal.WithExpectWorkerIdentityRegex(expectWorkerIdentityRegex),
	)

	return waiter.WaitContext(
		cmd.Context(),
		tc,
		waiter.WithTimeout(timeout),
		waiter.WithInterval(interval),
		waiter.WithInvertCheck(invertCheck),
		waiter.WithLogger(logger),
	)
}
