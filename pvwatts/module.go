package pvwatts

import (
	"encoding/json"
	"fmt"
	"gowatts/data"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	RetrieveSolarData(parameters *data.Parameters) (Output, error)
}

// New creates a new pvwatts struct
func New() API {
	apiKey := os.Getenv("PVWATTS_API_KEY")
	if apiKey == "" {
		fmt.Println("-------------- WARNING!! No PVWATTS API Api Key Set --------------------")
		fmt.Println("--- Using DEMO_KEY, this is limited to 50 requests per hour         ----")
		fmt.Println("--- Set the PVWATTS_API_KEY environment variable to use another key ----")
		fmt.Println("------------------------------------------------------------------------")
		apiKey = "DEMO_KEY"
	}
	return &pvWatts{apiKey}
}

type pvWatts struct {
	apiKey string
}

var httpGet = http.Get

// RetrieveSolarData sends a request to the API API and returns the result
func (p *pvWatts) RetrieveSolarData(parameters *data.Parameters) (Output, error) {
	var output Output

	// Do not continue if no location is set
	if parameters.Latitude == "" || parameters.Longitude == "" {
		return output, nil
	}

	// Send request to PVWatts API
	request := p.getPVWattsRequestURI(parameters)
	resp, err := httpGet(request)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	// Read the response of PVWatts
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return output, fmt.Errorf(string(body))
	}

	// Check if the response contains any errors
	err = json.Unmarshal(body, &output)
	if resp.StatusCode != http.StatusOK || len(output.Errors) > 0 || err != nil {
		return output, fmt.Errorf(string(body))
	}

	return output, nil
}

func (p *pvWatts) getPVWattsRequestURI(parameters *data.Parameters) string {
	dataset := getDataSet(parameters.Longitude, parameters.Latitude)
	baseURL := "https://developer.nrel.gov/api/pvwatts/v6.json"
	settings := fmt.Sprintf("api_key=%s&dataset=%s&timeframe=monthly&radius=0", p.apiKey, dataset)
	category := fmt.Sprintf("array_type=%s&module_type=%s", parameters.ArrayType, parameters.ModuleType)
	location := fmt.Sprintf("lat=%s&lon=%s", parameters.Latitude, parameters.Longitude)
	orientation := fmt.Sprintf("azimuth=%s&tilt=%s", parameters.Azimuth, parameters.Tilt)
	performance := fmt.Sprintf("system_capacity=%s&losses=%s", parameters.Capacity, parameters.Losses)

	return fmt.Sprintf("%s?%s&%s&%s&%s&%s", baseURL, settings, location, category, orientation, performance)
}

func getDataSet(longitude, latitude string) string {
	var lon, lat float64
	var err error
	intenationalDataset := "intl"
	nsrdbDataset := "nsrdb" // nsrdb dateset is available for India and America: https://nsrdb.nrel.gov/map.html
	if lon, err = strconv.ParseFloat(longitude, 64); err != nil {
		return intenationalDataset
	}

	if lat, err = strconv.ParseFloat(latitude, 64); err != nil {
		return intenationalDataset
	}

	// physical model
	if lat > -20 && lat < 60 && lon > -180 && lon < -20 {
		return nsrdbDataset
	}

	// india suny model
	if lat > 5 && lat < 38 && lon > 65 && lon < 93 {
		return nsrdbDataset
	}

	return intenationalDataset
}
