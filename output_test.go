package pvoutput

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeOutput(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/output/normal")
	require.NoError(t, err)

	output, err := decodeOutput(string(data))
	if assert.NoError(t, err) {
		date, _ := time.Parse("20060102", "20110327")
		assert.Equal(t, date, output.Date)
		assert.Equal(t, 4413, output.Generated)
		// Efficiency 0.460
		assert.Equal(t, 1234, output.Exported)
		assert.Equal(t, 21859, output.Consumption)
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
