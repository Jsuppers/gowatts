package server

import (
	"fmt"
	"gowatts/data"
	"gowatts/pvwatts"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server begins a server
type Server interface {
	Start()
}

// New creates a new http server
func New() Server {
	return &httpServer{pvwatts.New()}
}

// httpServer holds the information used in the template
type httpServer struct {
	pvwattsAPI pvwatts.API
}

// Start starts the server
func (server *httpServer) Start() {
	router := gin.Default()
	router.LoadHTMLGlob("resources/templates/*")
	router.Static("/static", "resources/static")
	router.GET("/", server.processRequest)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Error starting server ", err.Error())
	}
}

func (server *httpServer) processRequest(context *gin.Context) {
	parameters := getParameters(context)

	solarData, err := server.pvwattsAPI.RetrieveSolarData(&parameters)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	arrayTypes := createOptions(parameters.ArrayType, "Fixed - Open Rack", "Fixed - Roof Mounted", "1-Axis", "1-Axis Backtracking", "2-Axis")
	moduleTypes := createOptions(parameters.ModuleType, "Standard", "Premium", "Thin film")
	labels := createLabels(len(solarData.Data.DC))

	context.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Labels":      labels,
		"Zoom":        parameters.Zoom,
		"Tilt":        parameters.Tilt,
		"Losses":      parameters.Losses,
		"Azimuth":     parameters.Azimuth,
		"Latitude":    parameters.Latitude,
		"Longitude":   parameters.Longitude,
		"Capacity":    parameters.Capacity,
		"ArrayTypes":  arrayTypes,
		"ModuleTypes": moduleTypes,
		"Station":     solarData.Station,
		"Data":        solarData.Data,
	})
}

func getParameters(context *gin.Context) data.Parameters {
	var parameters data.Parameters
	parameters.Zoom = context.DefaultQuery("zoom", "2")
	parameters.Tilt = context.DefaultQuery("tilt", "40")
	parameters.Losses = context.DefaultQuery("losses", "10")
	parameters.Azimuth = context.DefaultQuery("azimuth", "0")
	parameters.Latitude = context.Query("latitude")
	parameters.Longitude = context.Query("longitude")
	parameters.Capacity = context.DefaultQuery("capacity", "1")
	parameters.ArrayType = context.DefaultQuery("arrayType", "0")
	parameters.ModuleType = context.DefaultQuery("moduleType", "0")
	return parameters
}
