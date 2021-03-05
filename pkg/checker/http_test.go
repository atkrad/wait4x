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
