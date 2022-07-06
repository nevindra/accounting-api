package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func IRRRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/IRR")
	ping.POST("/", GetIRR)
}

type IRR struct {
	Investment float64 `json:"investment"`
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

	result := irr.calculateIRR(irr.Investment, irr.Cashflows, irr.FirstRate/100, irr.SecondRate/100, irr.Baseline/100)
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

func (data IRR) calculateIRR(investment float64, cashflows []int, interestRate1 float64, interestRate2 float64, baseline float64) float64 {
	var (
		PV1, PV2  float64
		rateDif   float64
		PVInitial float64
	)

	PV1, _ = CalculatePV(cashflows, interestRate1)
	PV2, _ = CalculatePV(cashflows, interestRate2)
	rateDif = math.Abs(interestRate1 - interestRate2)
	PVDifference := math.Abs(PV1 - PV2)
	fmt.Println("PV Different:", PVDifference)
	if PV1 > PV2 {
		PVInitial = PV1 - investment
	} else {
		PVInitial = PV2 - investment
	}
	fmt.Println("Selisih Modal:", PVInitial)
	fmt.Println(PVDifference / PVInitial)
	IRR := baseline + (PVInitial/PVDifference)*rateDif
	return IRR * 100
}
