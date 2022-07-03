package routes

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func IRRRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/IRR")
	ping.POST("/", GetIRR)
}

type IRR struct {
	Investment int     `json:"investment"`
	Cashflows  []int   `json:"cashflows"`
	Baseline   float64 `json:"baseline"`
	FirstRate  float64 `json:"firstRate"`
	SecondRate float64 `json:"SecondRate"`
}

func GetIRR(c *gin.Context) {
	var (
		irr IRR
	)
	err := c.ShouldBindJSON(&irr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	result := calculateIRR(irr.Investment, irr.Cashflows, irr.FirstRate/100, irr.SecondRate/100)
	if result > irr.Baseline {
		c.JSON(http.StatusOK, gin.H{
			"positive": true,
			"IRR":      result,
			"message":  "The IRR is higher than the baseline",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"positive": false,
			"IRR":      result,
			"message":  "The IRR is lower than the baseline",
		})
	}
}

func calculateIRR(investment int, cashflows []int, interestRate1 float64, interestRate2 float64) float64 {
	var (
		NPV1, NPV2 float64
	)

	NPV1, _ = calculateNPV(investment, cashflows, interestRate1)
	NPV2, _ = calculateNPV(investment, cashflows, interestRate2)

	rateDif := interestRate1 - interestRate2
	NPVDif := math.Abs(NPV1 - NPV2)
	PVInitial := NPV1 - float64(investment)
	IRR := interestRate1 + (NPVDif/PVInitial)*rateDif
	return IRR * 100
}
