// Copyright 2019-2025 The Wait4X Authors
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

package cname

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"wait4x.dev/v2/checker"
)

const server = "wait4x.dev"

func TestCheckExistenceCNAME(t *testing.T) {
	d := New(server)
	assert.Nil(t, d.Check(context.Background()))
}

func TestCorrectCNAME(t *testing.T) {
	d := New(server, WithExpectedDomains([]string{"wait4x.dev"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestIncorrectCNAME(t *testing.T) {
	var expectedError *checker.ExpectedError
	d := New(server, WithExpectedDomains([]string{"something wrong"}))
	assert.ErrorAs(t, d.Check(context.Background()), &expectedError)
}

func TestCustomNSCorrectCNAME(t *testing.T) {
	d := New(server, WithNameServer("8.8.8.8:53"), WithExpectedDomains([]string{"wait4x.dev"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestRegexCorrectCNAME(t *testing.T) {
	d := New(server, WithExpectedDomains([]string{".*wait4.*"}))
	assert.Nil(t, d.Check(context.Background()))
}
