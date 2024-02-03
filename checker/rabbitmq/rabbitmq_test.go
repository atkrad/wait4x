// Copyright 2024 The Wait4X Authors
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

// Package rabbitmq provides RabbitMQ checker.
package rabbitmq

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"testing"
	"time"
	"wait4x.dev/v2/checker"
)

// RabbitMQSuite is a test suite for RabbitMQ checker
type RabbitMQSuite struct {
	suite.Suite
	container *rabbitmq.RabbitMQContainer
}

// SetupSuite starts a RabbitMQ container
func (s *RabbitMQSuite) SetupSuite() {
	var err error
	s.container, err = rabbitmq.RunContainer(context.Background(), testcontainers.WithImage("rabbitmq:3.12.12-management-alpine"))
	s.Require().NoError(err)
}

// TearDownSuite stops the RabbitMQ container
func (s *RabbitMQSuite) TearDownSuite() {
	err := s.container.Terminate(context.Background())
	s.Require().NoError(err)
}

// TestIdentity tests the identity of the RabbitMQ checker
func (s *RabbitMQSuite) TestIdentity() {
	chk := New("amqp://guest:guest@127.0.0.1:5672/vhost")
	identity, err := chk.Identity()

	s.Require().NoError(err)
	s.Assert().Equal("127.0.0.1:5672", identity)
}

// TestInvalidIdentity tests the invalid identity of the RabbitMQ checker
func (s *RabbitMQSuite) TestInvalidIdentity() {
	chk := New("127.0.0.1:5672")
	_, err := chk.Identity()

	s.Assert().ErrorContains(err, `can't retrieve the checker identity: parse "127.0.0.1:5672"`)
}

// TestValidConnection tests the valid connection of the RabbitMQ server
func (s *RabbitMQSuite) TestInvalidConnection() {
	var expectedError *checker.ExpectedError
	chk := New("amqp://user:pass@127.0.0.1:5672/vhost")

	s.Assert().ErrorAs(chk.Check(context.Background()), &expectedError)
}

// TestValidAddress tests the valid address of the RabbitMQ server
func (s *RabbitMQSuite) TestValidConnection() {
	ctx := context.Background()

	endpoint, err := s.container.AmqpURL(ctx)
	s.Require().NoError(err)

	chk := New(endpoint, WithTimeout(5*time.Second), WithInsecureSkipTLSVerify(true))
	s.Assert().Nil(chk.Check(ctx))
}

// TestRabbitMQ runs the RabbitMQ test suite
func TestRabbitMQ(t *testing.T) {
	suite.Run(t, new(RabbitMQSuite))
}
