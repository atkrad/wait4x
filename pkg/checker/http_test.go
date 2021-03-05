package checker

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/atkrad/wait4x/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestHttpInvalidAddress(t *testing.T) {
	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := NewHTTP("http://not-exists.tld", 0, "", time.Second*5)
	hc.SetLogger(logger)

	assert.Equal(t, false, hc.Check())
}

func TestHttpValidAddress(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := NewHTTP(ts.URL, 0, "", time.Second*5)
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check())
}

func TestHttpInvalidStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := NewHTTP(ts.URL, http.StatusCreated, "", time.Second*5)
	hc.SetLogger(logger)

	assert.Equal(t, false, hc.Check())
}

func TestHttpValidStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := NewHTTP(ts.URL, http.StatusOK, "", time.Second*5)
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check())
}

func TestHttpInvalidBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Wait4X"))
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := NewHTTP(ts.URL, 0, "FooBar", time.Second*5)
	hc.SetLogger(logger)

	assert.Equal(t, false, hc.Check())
}

func TestHttpValidBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Wait4X is the best CLI tools. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla volutpat arcu malesuada lacus vulputate feugiat. Etiam vitae sem quis ligula consequat euismod. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus fringilla sapien non lacus volutpat sollicitudin. Donec sollicitudin sit amet purus ac rutrum. Nam nunc orci, luctus a sagittis."))
	}))
	defer ts.Close()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	hc := NewHTTP(ts.URL, 0, "Wait4X.+best.+tools", time.Second*5)
	hc.SetLogger(logger)

	assert.Equal(t, true, hc.Check())
}
