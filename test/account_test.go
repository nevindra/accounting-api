package main

import (
	"basic/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func PerformRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPaybackPeriodRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	PaybackPeriodRoutes(r)
	w := PerformRequest(r, "POST", "/v1/paybackperiod", getPaybackPeriodJSON())
	assert.Equal(t, http.StatusOK, w.Code)
	mockResponse := `{"accumulatedCashflow":950,"cashBeforePeriod":600,"message":"The payback period is calculated","paybackPeriod":3.4285714285714284}`
	assert.Equal(t, mockResponse, w.Body.String())
}

func PaybackPeriodRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.POST("/paybackperiod", routes.CalculatepaybackPeriod)
}

func TestNPVRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NPVRoutes(r)
	w := PerformRequest(r, "POST", "/v1/npv", getNPVJSON())
	assert.Equal(t, http.StatusOK, w.Code)
	mockResponse := `{"NPV":-2830.0572044988803,"PV":147169.94279550112,"message":"The NPV is lower than 0","positive":false}`
	assert.Equal(t, mockResponse, w.Body.String())
}

func NPVRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.POST("/npv", routes.GetNPV)
}

func getNPVJSON() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`{
    "investment": 150000,
    "period":5, 
    "cashflows": [60000,50000,40000,35000,28000],
    "interestRate": 16
}`))
}

func TestIRRRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	IRRRoutes(r)
	w := PerformRequest(r, "POST", "/v1/irr", getIRRJSON())
	assert.Equal(t, http.StatusOK, w.Code)
	mockResponse := `{"IRR":15.15274875675059,"message":"The IRR is higher than the baseline","positive":true}`
	assert.Equal(t, mockResponse, w.Body.String())
}

func IRRRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.POST("/irr", routes.GetIRR)
}

func getIRRJSON() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`{
    "investment": 150000,
    "period":5, 
    "cashflows": [60000,50000,40000,35000,28000],
    "firstRate":16,
    "secondRate":10,
    "baseline":10
    }`))
}

func getPaybackPeriodJSON() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`{
		"investment": 750,
		"period": 15,
		"cashflows": [150,200,250,350,300]
	}`))
}
