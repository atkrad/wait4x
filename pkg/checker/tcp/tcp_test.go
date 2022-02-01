// Copyright 2021 Mohammad Abdolirad
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

package tcp

import (
	"context"
	"github.com/atkrad/wait4x/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"testing"
)

func TestTcpValidAddress(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")

	go func() {
		defer ln.Close()
		_, _ = ln.Accept()
	}()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	tc := New(ln.Addr().String())
	tc.SetLogger(logger)

	assert.Equal(t, nil, tc.Check(context.TODO()))
}

func TestTcpInvalidAddress(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")

	go func() {
		defer ln.Close()
		_, _ = ln.Accept()
	}()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	tc := New(ln.Addr().String() + "0")
	tc.SetLogger(logger)

	assert.Error(t, DialErr, tc.Check(context.TODO()))
}
