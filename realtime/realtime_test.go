package realtime

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestRealtime(t *testing.T) {
	f, err := os.Open("../data/zipkin.one")
	assert.NoError(t, err, "Error accessing test data")
	rdr := bufio.NewReader(f)
	s, err := NewRealtime(rdr)
	assert.NoError(t, err, "Error creating Realtime")
	assert.NoError(t, s.Run("http://localhost:9411/api/v1/spans","foobar"))
}

func TestStop(t *testing.T) {
	f, err := os.Open("../data/zipkin.medium")
	assert.NoError(t, err, "Error accessing test data")
	rdr := bufio.NewReader(f)
	s, err := NewRealtime(rdr)
	assert.NoError(t, err, "Error creating Realtime")
	go func() {
		assert.NoError(t, s.Run("http://localhost:9411/api/v1/spans", "foobar"))
	}()
	time.Sleep(1*time.Second)
	s.Stop()
}

func TestRunForever(t *testing.T) {
	f, err := os.Open("../data/zipkin.medium")
	assert.NoError(t, err, "Error accessing test data")
	rdr := bufio.NewReader(f)
	s, err := NewRealtime(rdr)
	assert.NoError(t, err, "Error creating Realtime")
	go func() {
		time.Sleep(5 * time.Second)
		s.Stop()
	}()
	assert.NoError(t, s.RunForever("http://localhost:9411/api/v1/spans","foobar"))
}