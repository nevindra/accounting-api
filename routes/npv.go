package routes

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

type NPV struct {
	Investment   int     `json:"investment"`
	Period       int     `json:"period"`
	Cashflows    []int   `json:"cashflows"`
	InterestRate float64 `json:"interestRate"`
}

func NPVRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/NPV")
	ping.POST("/", GetNPV)
}

func GetNPV(c *gin.Context) {
	var npv NPV
	err := c.ShouldBindJSON(&npv)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	result, _ := calculateNPV(npv.Investment, npv.Cashflows, npv.InterestRate/100)

	if result > 0 {
		c.JSON(http.StatusOK, gin.H{
			"positive": true,
			"NPV":      result,
			"message":  "The NPV is higher than 0",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"positive": false,
			"NPV":      result,
			"message":  "The NPV is lower than 0",
		})
	}

}

func calculateNPV(investment int, cashflows []int, interestRate float64) (float64, []float64) {

	var presentValue []float64
	var result float64

	// calculate present value for each cashflow
	for i := 0; i < len(cashflows); i++ {
		result = float64(cashflows[i]) * 1 / math.Pow(1+interestRate, float64(i+1))
		presentValue = append(presentValue, result)
	}

	// caluclate NPV
	var sum float64
	for i := 0; i < len(presentValue); i++ {
		sum += presentValue[i]
	}
	return sum - float64(investment), presentValue
}
