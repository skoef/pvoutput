package pvoutput

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeOutput(t *testing.T) {
	var result string
	var err error
	o := NewOutput()

	_, err = o.Encode()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Date is required")
	}

	newValidOutput := func() Output {
		o := NewOutput()
		o.Date, _ = time.Parse("20060102", "20200818")
		return o
	}

	// set valid baseline
	o = newValidOutput()
	result, err = o.Encode()
	assert.NoError(t, err)
	assert.Equal(t, "d=20200818", result)

	// check generated
	o = newValidOutput()
	o.Generated = 5
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&g=5", result)

	// check exported
	o = newValidOutput()
	o.Exported = 6
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&e=6", result)

	// check peakpower
	o = newValidOutput()
	o.PeakPower = 7
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&pp=7", result)

	// check peakpower
	o = newValidOutput()
	o.PeakTime, _ = time.Parse("1504", "1234")
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&pt=12%3A34", result)

	// check condition
	o = newValidOutput()
	o.Condition = "Sunny"
	result, _ = o.Encode()
	assert.Equal(t, "cd=Sunny&d=20200818", result)

	// check mintemp
	o = newValidOutput()
	o.MinTemp = 0.8
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&tm=0.8", result)

	// check maxtemp
	o = newValidOutput()
	o.MaxTemp = 0.9
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&tx=0.9", result)

	// check comments
	o = newValidOutput()
	o.Comments = "test 123"
	result, _ = o.Encode()
	assert.Equal(t, "cm=test+123&d=20200818", result)

	// check import peak
	o = newValidOutput()
	o.ImportPeak = 10
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&ip=10", result)

	// check import off-peak
	o = newValidOutput()
	o.ImportOffPeak = 11
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&io=11", result)

	// check import shoulder
	o = newValidOutput()
	o.ImportShoulder = 12
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&is=12", result)

	// check import high-shoulder
	o = newValidOutput()
	o.ImportHighShoulder = 13
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&ih=13", result)

	// check consumed
	o = newValidOutput()
	o.Consumed = 14
	result, _ = o.Encode()
	assert.Equal(t, "c=14&d=20200818", result)

	// check export peak
	o = newValidOutput()
	o.ExportPeak = 15
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&ep=15", result)

	// check export off-peak
	o = newValidOutput()
	o.ExportOffPeak = 16
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&eo=16", result)

	// check export shoulder
	o = newValidOutput()
	o.ExportShoulder = 17
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&es=17", result)

	// check export high-shoulder
	o = newValidOutput()
	o.ExportHighShoulder = 18
	result, _ = o.Encode()
	assert.Equal(t, "d=20200818&eh=18", result)
}

func TestDecodeOutput(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/output/normal")
	require.NoError(t, err)

	output, err := decodeOutput(string(data))
	if assert.NoError(t, err) {
		date, _ := time.Parse("20060102", "20110327")
		assert.Equal(t, date, output.Date)
		assert.Equal(t, 4413, output.Generated)
		assert.Equal(t, 0.46, output.Efficiency)
		assert.Equal(t, 1234, output.Exported)
		assert.Equal(t, 21859, output.Consumed)
		assert.Equal(t, 2070, output.PeakPower)
		ptime, _ := time.Parse("15:04", "11:00")
		assert.Equal(t, ptime, output.PeakTime)
		assert.Equal(t, "Showers", output.Condition)
		assert.Equal(t, -3.0, output.MinTemp)
		assert.Equal(t, 6.0, output.MaxTemp)
		assert.Equal(t, 4220, output.ImportPeak)
		assert.Equal(t, 7308, output.ImportOffPeak)
		assert.Equal(t, 2030, output.ImportShoulder)
		assert.Equal(t, 3888, output.ImportHighShoulder)
	}

	data, err = ioutil.ReadFile("testdata/output/timeofexport")
	require.NoError(t, err)

	output, err = decodeOutput(string(data))
	if assert.NoError(t, err) {
		assert.Equal(t, 3220, output.ExportPeak)
		assert.Equal(t, 6308, output.ExportOffPeak)
		assert.Equal(t, 1030, output.ExportShoulder)
		assert.Equal(t, 30, output.ExportHighShoulder)
	}

	data, err = ioutil.ReadFile("testdata/output/insolation")
	require.NoError(t, err)

	output, err = decodeOutput(string(data))
	if assert.NoError(t, err) {
		assert.Equal(t, 12910, output.Insolation)
	}
}

func TestEncodeBatchOutput(t *testing.T) {
	var b BatchOutput

	// expect error for empty batch
	b = BatchOutput{}
	_, err := b.Encode()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Empty")
	}

	// batches that are too big should throw an error as well
	b = make(BatchOutput, (BatchOutputMaxSize + 1))
	_, err = b.Encode()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Max")
	}

	// example of documentation
	// "Send three outputs in a single batch request"
	b = BatchOutput{NewOutput(), NewOutput(), NewOutput()}
	b[0].Date, _ = time.Parse("20060102", "20150101")
	b[0].Generated = 1239
	b[1].Date, _ = time.Parse("20060102", "20150102")
	b[1].Generated = 1523
	b[2].Date, _ = time.Parse("20060102", "20150103")
	b[2].Generated = 2190

	result, err := b.Encode()
	if assert.NoError(t, err) {
		assert.Equal(t, "data=20150101,1239;20150102,1523;20150103,2190", result)
	}

	// example of documentation
	// "Send a single status with Generation Energy 850Wh, Energy Used 1100Wh and Temperature 10.4C to 20.5C"
	b = BatchOutput{NewOutput()}
	b[0].Date, _ = time.Parse("20060102", "20150101")
	b[0].Generated = 850
	b[0].Consumed = 1100
	b[0].MinTemp = 10.4
	b[0].MaxTemp = 20.5

	result, err = b.Encode()
	if assert.NoError(t, err) {
		assert.Equal(t, "data=20150101,850,,1100,,,,10.4,20.5", result)
	}
}
