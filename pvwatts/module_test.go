package pvwatts

//go:generate mockgen -package=mocks -destination=./../mocks/io_mock.go io ReadCloser

import (
	"fmt"
	"gowatts/data"
	"gowatts/mocks"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  API
	}{
		{"test no api key", func() {}, &pvWatts{"DEMO_KEY"}},
		{"test sets api key", func() { os.Setenv("PVWATTS_API_KEY", "testKey") }, &pvWatts{"testKey"}},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			defer os.Unsetenv("PVWATTS_API_KEY")
			if got := New(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("New() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getDataSet(t *testing.T) {
	type args struct {
		longitude string
		latitude  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test unknown longitude", args{"fake longitude", "30"}, "intl"},
		{"test unknown latitude", args{"30", "fake latitude"}, "intl"},
		{"uses nsrdb at New York", args{"-73.935242", "40.730610"}, "nsrdb"},
		{"uses nsrdb at New Delhi", args{"77.216721", "28.644800"}, "nsrdb"},
		{"uses intl at Auckland", args{"174.763336", "-36.848461"}, "intl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSet(tt.args.longitude, tt.args.latitude); got != tt.want {
				t.Errorf("getDataSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pvWatts_getPVWattsRequestURI(t *testing.T) {
	tests := []struct {
		name       string
		parameters *data.Parameters
		want       string
	}{
		{"correctly sets parameters",
			&data.Parameters{
				Tilt:       "tilt",
				Losses:     "losses",
				Azimuth:    "azimuth",
				Latitude:   "latitude",
				Longitude:  "longitude",
				Capacity:   "capacity",
				ArrayType:  "arrayType",
				ModuleType: "moduleType",
				Zoom:       "zoom",
			},
			"https://developer.nrel.gov/api/pvwatts/v6.json?api_key=apiKey&dataset=intl&timeframe=monthly" +
				"&radius=0&lat=latitude&lon=longitude&array_type=arrayType&module_type=moduleType" +
				"&azimuth=azimuth&tilt=tilt&system_capacity=capacity&losses=losses"},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			testModule := pvWatts{apiKey: "apiKey"}
			if got := testModule.getPVWattsRequestURI(test.parameters); got != test.want {
				t.Errorf("pvWatts.getPVWattsRequestURI() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_pvWatts_RetrieveSolarData(t *testing.T) {
	testLatitude := "30"
	testLongitude := "30"

	type ioutilReadAllArgs struct {
		BodyString string
		RespError  error
	}
	type httpGetArgs struct {
		StatusCode int
		RespError  error
	}
	tests := []struct {
		name              string
		parameters        *data.Parameters
		httpGetArgs       httpGetArgs
		ioutilReadAllArgs ioutilReadAllArgs
		setup             func(mockBody *mocks.MockReadCloser)
		want              Output
		wantErr           bool
	}{
		{
			"no location does not return output",
			&data.Parameters{}, httpGetArgs{}, ioutilReadAllArgs{},
			func(mockBody *mocks.MockReadCloser) {}, Output{}, false,
		},
		{
			"bad request gives error",
			&data.Parameters{Latitude: testLatitude, Longitude: testLongitude},
			httpGetArgs{http.StatusBadRequest, fmt.Errorf("error")}, ioutilReadAllArgs{},
			func(mockBody *mocks.MockReadCloser) {}, Output{}, true,
		},
		{
			"error reading response returns error",
			&data.Parameters{Latitude: testLatitude, Longitude: testLongitude},
			httpGetArgs{http.StatusOK, nil},
			ioutilReadAllArgs{"body", fmt.Errorf("error reading")},
			func(mockBody *mocks.MockReadCloser) {
				mockBody.EXPECT().Close().Times(1)
			}, Output{}, true,
		},
		{
			"error marshalling request",
			&data.Parameters{Latitude: testLatitude, Longitude: testLongitude},
			httpGetArgs{http.StatusOK, nil},
			ioutilReadAllArgs{"{", nil},
			func(mockBody *mocks.MockReadCloser) {
				mockBody.EXPECT().Close().Times(1)
			}, Output{}, true,
		},
		{
			"success",
			&data.Parameters{Latitude: testLatitude, Longitude: testLongitude},
			httpGetArgs{http.StatusOK, nil},
			ioutilReadAllArgs{BodyString: "{\"station_info\":{\"lat\":10,\"lon\":20}}", RespError: nil},
			func(mockBody *mocks.MockReadCloser) {
				mockBody.EXPECT().Close().Times(1)
			}, Output{Station: Station{Latitude: 10.0, Longitude: 20.0}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := tt
			mockCtrl := gomock.NewController(t)
			mockBody := mocks.NewMockReadCloser(mockCtrl)
			test.setup(mockBody)

			testModule := pvWatts{apiKey: "apiKey"}
			httpGet = func(url string) (*http.Response, error) {
				return &http.Response{StatusCode: test.httpGetArgs.StatusCode, Body: mockBody}, test.httpGetArgs.RespError
			}
			ioutilReadAll = func(r io.Reader) ([]byte, error) {
				return []byte(test.ioutilReadAllArgs.BodyString), test.ioutilReadAllArgs.RespError
			}
			defer func() {
				httpGet = http.Get
				ioutilReadAll = ioutil.ReadAll
			}()

			got, err := testModule.RetrieveSolarData(test.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("pvWatts.RetrieveSolarData() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("pvWatts.RetrieveSolarData() = %v, want %v", got, test.want)
			}
		})
	}
}
