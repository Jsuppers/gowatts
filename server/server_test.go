package server

import (
	"fmt"
	"gowatts/data"
	"gowatts/pvwatts"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Server
	}{
		{"creates new http server", &httpServer{pvwatts.New()}},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("New() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestHttpResponse(t *testing.T) {
	tests := []struct {
		name           string
		API            pvwatts.API
		expectedStatus int
	}{
		{"correctly sends response", &mockSuccessAPI{}, 200},
		{"correctly sends error if pvwatts API returns error", &mockFailAPI{}, 400},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			mockAPI := test.API
			server := &httpServer{mockAPI}

			templatesPath = "../resources/templates/*"
			router := server.setupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)

			router.ServeHTTP(w, req)

			if !reflect.DeepEqual(w.Code, test.expectedStatus) {
				t.Errorf("New() = %v, want %v", w.Code, test.expectedStatus)
			}
		})
	}
}

type mockSuccessAPI struct{}

func (p *mockSuccessAPI) RetrieveSolarData(parameters *data.Parameters) (pvwatts.Output, error) {
	return pvwatts.Output{}, nil
}

type mockFailAPI struct{}

func (p *mockFailAPI) RetrieveSolarData(parameters *data.Parameters) (pvwatts.Output, error) {
	return pvwatts.Output{}, fmt.Errorf("cannot reach pvwatts API")
}
