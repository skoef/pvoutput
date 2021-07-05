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
	outputUnsetInt    = -1
	outputUnsetFloat  = -1.0
	outputUnsetString = "__unset__"
	// order of keys in batch output
	// as described on https://pvoutput.org/help.html#api-addbatchoutput
	outputBatchKeys = []string{
		"d",  // date
		"g",  // generated
		"e",  // exported
		"c",  // consumed
		"pp", // peak power
		"pt", // peak time
		"cd", // condition
		"tm", // min temperature
		"tx", // max temperature
		"cm", // comments
		"ip", // import peak
		"io", // import off-peak
		"is", // import shoulder
	}
)

const (
	// BatchOutputMaxSize determines the maximum batch size
	// this is 30 according to PVOutput's docs
	BatchOutputMaxSize = 30
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
	Consumed           int // watt hours
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
		Consumed:           outputUnsetInt,
		ExportPeak:         outputUnsetInt,
		ExportOffPeak:      outputUnsetInt,
		ExportShoulder:     outputUnsetInt,
		ExportHighShoulder: outputUnsetInt,
	}
}

func (o Output) encode() (url.Values, error) {
	data := url.Values{}
	if o.Date.IsZero() {
		return nil, errors.New("date is required on Output")
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
	if o.Consumed != outputUnsetInt {
		data.Set("c", fmt.Sprintf("%d", o.Consumed))
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

	return data, nil
}

// Encode returns API string for this object
func (o Output) Encode() (string, error) {
	data, err := o.encode()
	if err != nil {
		return "", err
	}

	return data.Encode(), nil
}

func decodeOutput(input string) (op Output, err error) {
	fields := strings.Split(strings.TrimSpace(input), ",")
	if len(fields) < 14 {
		return
	}
	// parse Date field from fields[0]
	op.Date, err = time.Parse("20060102", fields[0])
	if err != nil {
		return
	}
	// parse Generated field from fields[1]
	op.Generated, err = strconv.Atoi(fields[1])
	if err != nil {
		return
	}

	// parse Efficiency field from fields[8]
	op.Efficiency, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return
	}

	// parse Exported field from fields[3]
	op.Exported, err = strconv.Atoi(fields[3])
	if err != nil {
		return
	}

	// parse Consumed field from fields[4]
	op.Consumed, err = strconv.Atoi(fields[4])
	if err != nil {
		return
	}

	// parse PeakPower field from fields[5]
	op.PeakPower, err = strconv.Atoi(fields[5])
	if err != nil {
		return
	}

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
	op.ImportPeak, err = strconv.Atoi(fields[10])
	if err != nil {
		return
	}

	// parse ImportOffPeak field from fields[11]
	op.ImportOffPeak, err = strconv.Atoi(fields[11])
	if err != nil {
		return
	}

	// parse ImportShoulder field from fields[12]
	op.ImportShoulder, err = strconv.Atoi(fields[12])
	if err != nil {
		return
	}

	// parse ImportHighShoulder field from fields[13]
	op.ImportHighShoulder, err = strconv.Atoi(fields[13])
	if err != nil {
		return
	}

	if len(fields) < 18 {
		return
	}

	// parse ExportPeak field from fields[14]
	op.ExportPeak, err = strconv.Atoi(fields[14])
	if err != nil {
		return
	}

	// parse ExportOffPeak field from fields[15]
	op.ExportOffPeak, err = strconv.Atoi(fields[15])
	if err != nil {
		return
	}

	// parse ExportShoulder field from fields[16]
	op.ExportShoulder, err = strconv.Atoi(fields[16])
	if err != nil {
		return
	}

	// parse ExportHighShoulder field from fields[17]
	op.ExportHighShoulder, err = strconv.Atoi(fields[17])
	if err != nil {
		return
	}

	if len(fields) < 19 {
		return
	}

	// parse Insolation field from fields[18]
	op.Insolation, err = strconv.Atoi(fields[18])
	if err != nil {
		return
	}

	return
}

// BatchOutput is a convenience type for a slice of Outputs
type BatchOutput []Output

// Encode returns API string for this object
func (b BatchOutput) Encode() (string, error) {
	if len(b) == 0 {
		return "", errors.New("empty batch")
	}

	if len(b) > BatchOutputMaxSize {
		return "", fmt.Errorf("max batch size is %d", BatchOutputMaxSize)
	}

	items := []string{}

	for _, o := range b {
		fields := []string{}
		enc, err := o.encode()
		if err != nil {
			return "", err
		}

		for _, key := range outputBatchKeys {
			fields = append(fields, enc.Get(key))
		}

		line := strings.Join(fields, ",")
		items = append(items, strings.TrimRight(line, ","))
	}

	return fmt.Sprintf("data=%s", strings.Join(items, ";")), nil
}
