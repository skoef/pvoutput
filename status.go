package pvoutput

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// StatusCumulative is a flag to tell if and how a status update has cumulative Wh values
type StatusCumulative int

const (
	// StatusCumulativeAll all Wh values are lifetime energy values
	StatusCumulativeAll StatusCumulative = 1
	// StatusCumulativeGenerating generation Wh values are lifetime energy values
	StatusCumulativeGenerating StatusCumulative = 2
	// StatusCumulativeConsuming consumption Wh values are lifetime energy values
	StatusCumulativeConsuming StatusCumulative = 3
)

// Status represents the data structure for a PV status update as described
// on https://pvoutput.org/help.html#api-addstatus
type Status struct {
	DateTime    time.Time
	Generated   int     // watt hours
	Generating  int     // watts
	Consumed    int     // watt hours
	Consuming   int     // watts
	Output      float64 // kW / kW ratio
	Temperature float64 // celsius
	Voltage     float64 // volts
	Cumulative  StatusCumulative
}

// NewStatus initialises and returns a new Status
// the reason why we set everything to "unset" values
// is to detect the difference between an unset field
// or a deliberately set to default values of
// 0 and 0.0 for int and float respectively
// this helps during encoding the status update to API POST
// body
func NewStatus() Status {
	return Status{
		Generated:   outputUnsetInt,
		Generating:  outputUnsetInt,
		Consumed:    outputUnsetInt,
		Consuming:   outputUnsetInt,
		Temperature: outputUnsetFloat,
		Voltage:     outputUnsetFloat,
		Cumulative:  StatusCumulative(outputUnsetInt),
	}
}

// Encode returns API string for this object
func (s Status) Encode() (string, error) {
	data := url.Values{}
	if s.DateTime.IsZero() {
		return "", errors.New("DateTime is required on Status")
	}

	data.Set("d", s.DateTime.Format("20060102"))
	data.Set("t", s.DateTime.Format("15:04"))

	if s.Generated != outputUnsetInt {
		data.Set("v1", fmt.Sprintf("%d", s.Generated))
	}
	if s.Generating != outputUnsetInt {
		data.Set("v2", fmt.Sprintf("%d", s.Generating))
	}
	if s.Consumed != outputUnsetInt {
		data.Set("v3", fmt.Sprintf("%d", s.Consumed))
	}
	if s.Consuming != outputUnsetInt {
		data.Set("v4", fmt.Sprintf("%d", s.Consuming))
	}
	if s.Temperature != outputUnsetFloat {
		data.Set("v5", fmt.Sprintf("%0.1f", s.Temperature))
	}
	if s.Voltage != outputUnsetFloat {
		data.Set("v6", fmt.Sprintf("%0.1f", s.Voltage))
	}
	if int(s.Cumulative) != outputUnsetInt {
		data.Set("c1", fmt.Sprintf("%d", s.Cumulative))
	}

	return data.Encode(), nil
}

func decodeStatus(input string) (s Status, err error) {
	var plh int64
	fields := strings.Split(strings.TrimSpace(input), ",")

	if len(fields) >= 9 {
		// parse DateTime field from fields[0]+fields[1]
		s.DateTime, err = time.Parse("20060102-15:04", fmt.Sprintf("%s-%s", fields[0], fields[1]))
		if err != nil {
			return
		}
		// parse Generated field from fields[2]
		plh, err = strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return
		}
		s.Generated = int(plh)
		// parse Generating field from fields[3]
		plh, err = strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			return
		}
		s.Generating = int(plh)
		// parse Consumed field from fields[4]
		plh, err = strconv.ParseInt(fields[4], 10, 64)
		if err != nil {
			return
		}
		s.Consumed = int(plh)
		// parse Consuming field from fields[5]
		plh, err = strconv.ParseInt(fields[5], 10, 64)
		if err != nil {
			return
		}
		s.Consuming = int(plh)
		// parse Output field from fields[6]
		s.Output, err = strconv.ParseFloat(fields[6], 64)
		if err != nil {
			return
		}
		// parse Temperature field from fields[7]
		s.Temperature, err = strconv.ParseFloat(fields[7], 64)
		if err != nil {
			return
		}
		// parse Voltage field from fields[8]
		s.Voltage, err = strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return
		}
	}

	return
}
