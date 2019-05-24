package reader

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSmall(t *testing.T) {
	f, err := os.Open("../data/zipkin.one")
	assert.NoError(t, err, "Error accessing test data")
	rdr := bufio.NewReader(f)
	dr := NewDumpReader(rdr)

	i := 0
	for {
		r, err := dr.Next()
		assert.NoError(t, err)
		if r == nil {
			break
		}
		assert.NotNil(t, r, "Returned span was nil")
		assert.NotEqual(t, 0, len(r.Span), "Span was empty")

		// Spot checks of the data
		for _, s := range r.Span {
			assert.NotEmpty(t, s.TraceID, "TraceID missing")
			assert.NotEmpty(t, s.Name, "Name missing")
		}
		i++
	}
	fmt.Println(i)
}


