package pvwatts

import (
	"encoding/json"
	"fmt"
	"gowatts/data"
	"io/ioutil"
	"net/http"
	"os"
)

// Output holds the output from the API API call
type Output struct {
	Station Station  `json:"station_info"`
	Errors  []string `json:"errors"`
	Data    Data     `json:"outputs"`
}

// Station holds information about the weather station which was used
type Station struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

// Data holds various solar energy data for different time periods
type Data struct {
	AC  []float64 `json:"ac_monthly"`
	DC  []float64 `json:"dc_monthly"`
	POA []float64 `json:"poa_monthly"`
}

// API interface allows interaction with the API API
type API interface {
	RetrieveSolarData(parameters data.Parameters) (Output, error)
}

// New creates a new pvwatts struct
func New() API {
	apiKey := os.Getenv("PVWATTS_API_KEY")
	if apiKey == "" {
		fmt.Println("---------- WARNING!! -----------------")
		fmt.Println("\n---------- No API Api Key Set --------")
		fmt.Println("Using DEMO_KEY, this is limited to 50 requests per hour")
		fmt.Println("To set another key change the PVWATTS_API_KEY environment variable")
		fmt.Println("------------------------------------------")
		apiKey = "DEMO_KEY"
	}
	return &pvWatts{apiKey}
}

type pvWatts struct {
	apiKey string
}

// RetrieveSolarData sends a request to the API API and returns the result
func (p *pvWatts) RetrieveSolarData(parameters data.Parameters) (Output, error) {
	var output Output

	if parameters.Latitude == "" || parameters.Longitude == "" {
		return output, nil
	}

	timeframe := "monthly" //hourly
	dataset := "intl"      // TODO check is within usa

	baseURL := fmt.Sprintf("https://developer.nrel.gov/api/pvwatts/v6.json?api_key=%s&dataset=%s&timeframe=%s&radius=0&", p.apiKey, dataset, timeframe)
	param := fmt.Sprintf("lat=%s&lon=%s&system_capacity=%s&azimuth=%s&tilt=%s&array_type=%s&module_type=%s&losses=%s", parameters.Latitude, parameters.Longitude, parameters.Capacity, parameters.Azimuth, parameters.Tilt, parameters.ArrayType, parameters.ModuleType, parameters.Losses)
	request := fmt.Sprintf("%s%s", baseURL, param)

	fmt.Println("Sending Request")
	fmt.Println(request)

	resp, err := http.Get(request)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &output)

	if len(output.Errors) > 0 {
		return output, fmt.Errorf(output.Errors[0])
	}

	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf(string(body))
	}
	fmt.Println(string(body))

	return output, nil

}
