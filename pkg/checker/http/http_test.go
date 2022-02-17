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

package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atkrad/wait4x/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestHttpInvalidAddress(t *testing.T) {
	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := New("http://not-exists.tld")
	hc.SetLogger(logger)

	assert.Equal(t, false, hc.Check(context.TODO()))
}

func TestHttpValidAddress(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := New(ts.URL)
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check(context.TODO()))
}

func TestHttpInvalidStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := New(ts.URL, WithExpectStatusCode(http.StatusCreated))
	hc.SetLogger(logger)

	assert.Equal(t, false, hc.Check(context.TODO()))
}

func TestHttpValidStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := New(ts.URL, WithExpectStatusCode(http.StatusOK))
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check(context.TODO()))
}

func TestHttpInvalidBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Wait4X"))
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := New(ts.URL, WithExpectBody("FooBar"))
	hc.SetLogger(logger)

	assert.Equal(t, false, hc.Check(context.TODO()))
}

func TestHttpValidBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Wait4X is the best CLI tools. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla volutpat arcu malesuada lacus vulputate feugiat. Etiam vitae sem quis ligula consequat euismod. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus fringilla sapien non lacus volutpat sollicitudin. Donec sollicitudin sit amet purus ac rutrum. Nam nunc orci, luctus a sagittis."))
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := New(ts.URL, WithExpectBody("Wait4X.+best.+tools"))
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check(context.TODO()))
}

func TestHttpValidHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Test-Header", "test-value")
		w.Header().Add("Test-Header-New", "test-value-new")
		w.Header().Add("Authorization", "Token 1234")
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := New(ts.URL, WithExpectHeader("Test-Header"))
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check(context.TODO()))

	// Regex.
	hc = New(ts.URL, WithExpectHeader("Test-Header=test-.+"))
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check(context.TODO()))

	//hc = New(ts.URL, WithExpectHeader("Authorization=^Token\\s.+"))
	hc = New(ts.URL, WithExpectHeader("Authorization=^Token\\s.+"))
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check(context.TODO()))

	// Key value.
	hc = New(ts.URL, WithExpectHeader("Test-Header=test-value"))
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check(context.TODO()))
}

func TestHttpInvalidHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Test-Header", "test-value")
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := New(ts.URL, WithExpectHeader("Test-Header-New"))
	hc.SetLogger(logger)

	assert.Equal(t, false, hc.Check(context.TODO()))

	hc = New(ts.URL, WithExpectHeader("Test-.+=test-value"))
	hc.SetLogger(logger)

	assert.Equal(t, false, hc.Check(context.TODO()))
}
