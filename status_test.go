package pvoutput

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusEncode(t *testing.T) {
	var result string
	var err error
	s := NewStatus()

	_, err = s.Encode()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "DateTime is required")
	}

	newValidStatus := func() Status {
		s := NewStatus()
		s.DateTime, _ = time.Parse("20060102T15:04", "20200818T12:34")
		return s
	}

	// set valid baseline
	s = newValidStatus()
	result, err = s.Encode()
	assert.NoError(t, err)
	assert.Equal(t, "d=20200818&t=12%3A34", result)

	// test generated
	s = newValidStatus()
	s.Generated = 1
	result, _ = s.Encode()
	assert.Equal(t, "d=20200818&t=12%3A34&v1=1", result)

	// test generating
	s = newValidStatus()
	s.Generating = 2
	result, _ = s.Encode()
	assert.Equal(t, "d=20200818&t=12%3A34&v2=2", result)

	// test consumed
	s = newValidStatus()
	s.Consumed = 3
	result, _ = s.Encode()
	assert.Equal(t, "d=20200818&t=12%3A34&v3=3", result)

	// test consuming
	s = newValidStatus()
	s.Consuming = 4
	result, _ = s.Encode()
	assert.Equal(t, "d=20200818&t=12%3A34&v4=4", result)

	// test temperature
	s = newValidStatus()
	s.Temperature = 5.0213
	result, _ = s.Encode()
	assert.Equal(t, "d=20200818&t=12%3A34&v5=5.0", result)

	// test voltage
	s = newValidStatus()
	s.Voltage = 6.1234
	result, _ = s.Encode()
	assert.Equal(t, "d=20200818&t=12%3A34&v6=6.1", result)

	// test cumulative
	s = newValidStatus()
	s.Cumulative = StatusCumulativeConsuming
	result, _ = s.Encode()
	assert.Equal(t, "c1=3&d=20200818&t=12%3A34", result)
}

func TestDecodeStatus(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/status/normal")
	require.NoError(t, err)

	status, err := decodeStatus(string(data))
	if assert.NoError(t, err) {
		dtime, _ := time.Parse("200601021504", "201011071830")
		assert.Equal(t, dtime, status.DateTime)
		assert.Equal(t, 12936, status.Generated)
		assert.Equal(t, 202, status.Generating)
		assert.Equal(t, 19832, status.Consumed)
		assert.Equal(t, 459, status.Consuming)
		assert.Equal(t, 5.28, status.Output)
		assert.Equal(t, 15.3, status.Temperature)
		assert.Equal(t, 240.1, status.Voltage)
	}
}
