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

func TestEncodeBatchStatus(t *testing.T) {
	var b BatchStatus

	// expect error for empty batch
	b = BatchStatus{}
	_, err := b.Encode()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "empty")
	}

	// batches that are too big should throw an error as well
	b = make(BatchStatus, (BatchStatusMaxSize + 1))
	_, err = b.Encode()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "max")
	}

	// example of documentation
	// "Send three statuses from 10:00AM to 10:10AM in a single batch request"
	b = BatchStatus{NewStatus(), NewStatus(), NewStatus()}
	b[0].DateTime, _ = time.Parse("200601021504", "201101121000")
	b[0].Generated = 705
	b[0].Generating = 1029
	b[1].DateTime, _ = time.Parse("200601021504", "201101121005")
	b[1].Generated = 775
	b[1].Generating = 1320
	b[2].DateTime, _ = time.Parse("200601021504", "201101121010")
	b[2].Generated = 800
	b[2].Generating = 800

	result, err := b.Encode()
	if assert.NoError(t, err) {
		assert.Equal(t, "data=20110112,10:00,705,1029;20110112,10:05,775,1320;20110112,10:10,800,800", result)
	}

	// example of documentation
	// "Send a single status with Generation Energy 850Wh, Generation Power 1109W, Temperature 23.1C and Voltage 240V"
	b = BatchStatus{NewStatus()}
	b[0].DateTime, _ = time.Parse("200601021504", "201101121015")
	b[0].Generated = 850
	b[0].Generating = 1109
	b[0].Temperature = 23.1
	b[0].Voltage = 240

	result, err = b.Encode()
	if assert.NoError(t, err) {
		assert.Equal(t, "data=20110112,10:15,850,1109,,,23.1,240.0", result)
	}

	// example of documentation
	// "Send a single status with Consumption Energy 2000Wh, Consumption Power 210W"
	// data=20110112,4:15,,,2000,210
	b = BatchStatus{NewStatus()}
	b[0].DateTime, _ = time.Parse("200601021504", "201101120415")
	b[0].Consumed = 2000
	b[0].Consuming = 210

	result, err = b.Encode()
	if assert.NoError(t, err) {
		assert.Equal(t, "data=20110112,04:15,,,2000,210", result)
	}
}
