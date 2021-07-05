package pvoutput

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	a := NewAPI("foo", "bar", true)
	assert.Equal(t, "foo", a.Key)
	assert.Equal(t, "bar", a.SystemID)
	assert.True(t, a.donating)
}

func TestAPIAddBatchOutput(t *testing.T) {
	// create batch output of more than 1 item when not in donating mode
	b := make(BatchOutput, (BatchOutputMaxSize + 1))
	a := API{}
	err := a.AddBatchOutput(b)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Sprintf("max batch size is %d", BatchOutputMaxSize), err.Error())
	}

	// in donating mode, add more outputs and trigger same error
	a.donating = true
	b = make(BatchOutput, (BatchOutputMaxSizeDonating + 1))
	err = a.AddBatchOutput(b)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Sprintf("max batch size is %d", BatchOutputMaxSizeDonating), err.Error())
	}
}
