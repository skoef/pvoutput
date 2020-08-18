package pvoutput

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	outputUnsetInt    int     = -1
	outputUnsetFloat  float64 = -1.0
	outputUnsetString string  = "__unset__"
)

// Output represents the data structure for a PV Output as described
// on https://pvoutput.org/help.html#api-addoutput
type Output struct {
	Date               time.Time
	Generated          int     // watt hours
	Efficiency         float64 // ratio
	Exported           int     // watt hours
	PeakPower          int     // watts
	PeakTime           time.Time
	Condition          string
	MinTemp            float64 // degrees celsius
	MaxTemp            float64 // degrees celsius
	Comments           string
	ImportPeak         int // watt hours
	ImportOffPeak      int // watt hours
	ImportShoulder     int // watt hours
	ImportHighShoulder int // watt hours
	Consumption        int // watt hours
	ExportPeak         int // watt hours
	ExportOffPeak      int // watt hours
	ExportShoulder     int // watt hours
	ExportHighShoulder int // watt hours
	Insolation         int // watt hours
}

// NewOutput initialises and returns a new Output
// the reason why we set everything to "unset" values
// is to detect the difference between an unset field
// or a deliberately set to default values of
// 0 and 0.0 for int and float respectively
// this helps during encoding the output to API POST
// body
func NewOutput() Output {
	return Output{
		Generated:          outputUnsetInt,
		Exported:           outputUnsetInt,
		PeakPower:          outputUnsetInt,
		Condition:          outputUnsetString,
		MinTemp:            outputUnsetFloat,
		MaxTemp:            outputUnsetFloat,
		Comments:           outputUnsetString,
		ImportPeak:         outputUnsetInt,
		ImportOffPeak:      outputUnsetInt,
		ImportShoulder:     outputUnsetInt,
		ImportHighShoulder: outputUnsetInt,
		Consumption:        outputUnsetInt,
		ExportPeak:         outputUnsetInt,
		ExportOffPeak:      outputUnsetInt,
		ExportShoulder:     outputUnsetInt,
		ExportHighShoulder: outputUnsetInt,
	}
}

// Encode returns API string for this object
func (o Output) Encode() (string, error) {
	data := url.Values{}
	if o.Date.IsZero() {
		return "", errors.New("Date is required on Output")
	}

	data.Set("d", o.Date.Format("20060102"))
	if o.Generated != outputUnsetInt {
		data.Set("g", fmt.Sprintf("%d", o.Generated))
	}
	if o.Exported != outputUnsetInt {
		data.Set("e", fmt.Sprintf("%d", o.Exported))
	}
	if o.PeakPower != outputUnsetInt {
		data.Set("pp", fmt.Sprintf("%d", o.PeakPower))
	}
	if !o.PeakTime.IsZero() {
		data.Set("pt", o.PeakTime.Format("15:04"))
	}
	if o.Condition != outputUnsetString {
		data.Set("cd", o.Condition)
	}
	if o.MinTemp != outputUnsetFloat {
		data.Set("tm", fmt.Sprintf("%0.1f", o.MinTemp))
	}
	if o.MaxTemp != outputUnsetFloat {
		data.Set("tx", fmt.Sprintf("%0.1f", o.MaxTemp))
	}
	if o.Comments != outputUnsetString {
		data.Set("cm", o.Comments)
	}
	if o.ImportPeak != outputUnsetInt {
		data.Set("ip", fmt.Sprintf("%d", o.ImportPeak))
	}
	if o.ImportOffPeak != outputUnsetInt {
		data.Set("io", fmt.Sprintf("%d", o.ImportOffPeak))
	}
	if o.ImportShoulder != outputUnsetInt {
		data.Set("is", fmt.Sprintf("%d", o.ImportShoulder))
	}
	if o.ImportHighShoulder != outputUnsetInt {
		data.Set("ih", fmt.Sprintf("%d", o.ImportHighShoulder))
	}
	if o.Consumption != outputUnsetInt {
		data.Set("c", fmt.Sprintf("%d", o.Consumption))
	}
	if o.ExportPeak != outputUnsetInt {
		data.Set("ep", fmt.Sprintf("%d", o.ExportPeak))
	}
	if o.ExportOffPeak != outputUnsetInt {
		data.Set("eo", fmt.Sprintf("%d", o.ExportOffPeak))
	}
	if o.ExportShoulder != outputUnsetInt {
		data.Set("es", fmt.Sprintf("%d", o.ExportShoulder))
	}
	if o.ExportHighShoulder != outputUnsetInt {
		data.Set("eh", fmt.Sprintf("%d", o.ExportHighShoulder))
	}

	return data.Encode(), nil
}

func decodeOutput(input string) (op Output, err error) {
	var plh int64
	fields := strings.Split(strings.TrimSpace(input), ",")
	if len(fields) >= 14 {
		// parse Date field from fields[0]
		op.Date, err = time.Parse("20060102", fields[0])
		if err != nil {
			return
		}
		// parse Generated field from fields[1]
		plh, err = strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return
		}
		op.Generated = int(plh)
		// parse Efficiency field from fields[8]
		op.Efficiency, err = strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return
		}
		// parse Exported field from fields[3]
		plh, err = strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			return
		}
		op.Exported = int(plh)
		// parse Consumed field from fields[4]
		plh, err = strconv.ParseInt(fields[4], 10, 64)
		if err != nil {
			return
		}
		op.Consumption = int(plh)
		// parse PeakPower field from fields[5]
		plh, err = strconv.ParseInt(fields[5], 10, 64)
		if err != nil {
			return
		}
		op.PeakPower = int(plh)

		// parse PeakTime field from fields[6]
		op.PeakTime, err = time.Parse("15:04", fields[6])
		if err != nil {
			return
		}
		// get Condition field from fields[7]
		op.Condition = fields[7]
		// parse MinTemp field from fields[8]
		op.MinTemp, err = strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return
		}
		// parse MaxTemp field from fields[9]
		op.MaxTemp, err = strconv.ParseFloat(fields[9], 64)
		if err != nil {
			return
		}
		// parse ImportPeak field from fields[10]
		plh, err = strconv.ParseInt(fields[10], 10, 64)
		if err != nil {
			return
		}
		op.ImportPeak = int(plh)
		// parse ImportOffPeak field from fields[11]
		plh, err = strconv.ParseInt(fields[11], 10, 64)
		if err != nil {
			return
		}
		op.ImportOffPeak = int(plh)
		// parse ImportShoulder field from fields[12]
		plh, err = strconv.ParseInt(fields[12], 10, 64)
		if err != nil {
			return
		}
		op.ImportShoulder = int(plh)
		// parse ImportHighShoulder field from fields[13]
		plh, err = strconv.ParseInt(fields[13], 10, 64)
		if err != nil {
			return
		}
		op.ImportHighShoulder = int(plh)
	}

	if len(fields) >= 18 {
		// parse ExportPeak field from fields[14]
		plh, err = strconv.ParseInt(fields[14], 10, 64)
		if err != nil {
			return
		}
		op.ExportPeak = int(plh)
		// parse ExportOffPeak field from fields[15]
		plh, err = strconv.ParseInt(fields[15], 10, 64)
		if err != nil {
			return
		}
		op.ExportOffPeak = int(plh)
		// parse ExportShoulder field from fields[16]
		plh, err = strconv.ParseInt(fields[16], 10, 64)
		if err != nil {
			return
		}
		op.ExportShoulder = int(plh)
		// parse ExportHighShoulder field from fields[17]
		plh, err = strconv.ParseInt(fields[17], 10, 64)
		if err != nil {
			return
		}
		op.ExportHighShoulder = int(plh)
	}

	if len(fields) >= 19 {
		// parse Insolation field from fields[18]
		plh, err = strconv.ParseInt(fields[18], 10, 64)
		if err != nil {
			return
		}
		op.Insolation = int(plh)
	}

	return
}
