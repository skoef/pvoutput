[![Go Report Card](https://goreportcard.com/badge/github.com/skoef/pvoutput)](https://goreportcard.com/report/github.com/skoef/pvoutput) [![Documentation](https://godoc.org/github.com/skoef/pvoutput?status.svg)](http://godoc.org/github.com/skoef/pvoutput)

# Golang API Client for PVOutput.org

This is a golang library to interact with the [PVOutput](https://pvoutput.org) API.

## Example usage:
```golang
package main

import (
	"time"

	"github.com/skoef/pvoutput"
)

func main() {
    // create an API client per system you want to manage
    api := pvoutput.NewAPI("XXX", "12345", false)

    // get PV generation data from your solar inverter or
    // power consumption data from your utility meter
    //
    // ...

    // format data in pvoutput struct
    output := pvoutput.NewOutput()
    output.Date = time.Now()
    output.Generated = 4567 // 4.5 kWh
    output.Consumed = 7124 // 7.1 kWh

    // write data to API
    err := api.AddOutput(output)
    if err != nil {
        panic(err)
    }
}
```
