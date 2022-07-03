package routes

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func IRRRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/IRR")
	ping.POST("/", getIRR)
}

type IRR struct {
	Investment int     `json:"investment"`
	Cashflows  []int   `json:"cashflows"`
	FirstRate  float64 `json:"firstRate"`
	SecondRate float64 `json:"SecondRate"`
}

func getIRR(c *gin.Context) {
	var (
		irr IRR
	)
	c.ShouldBindJSON(&irr)

	result := calculateIRR(irr.Investment, irr.Cashflows, irr.FirstRate/100, irr.SecondRate/100)
	if result > 10 {
		c.JSON(http.StatusOK, gin.H{
			"status": "IRR is positve",
			"IRR":    result,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "IRR is negative",
			"IRR":    result,
		})
	}
}

func calculateIRR(investment int, cashflows []int, interestRate1 float64, interestRate2 float64) float64 {
	var (
		NPV1, NPV2 float64
	)

	NPV1, _ = calculateNPV(investment, cashflows, interestRate1)
	NPV2, _ = calculateNPV(investment, cashflows, interestRate2)
	//fmt.Println("NPV1: ", NPV1)
	//fmt.Println("NPV2: ", NPV2)

	rateDif := interestRate1 - interestRate2
	//fmt.Println(rateDif)
	NPVDif := math.Abs(NPV1 - NPV2)
	//fmt.Println(NPVDif)
	PVInitial := NPV1 - float64(investment)
	//fmt.Println(PVInitial)
	IRR := interestRate1 + (NPVDif/PVInitial)*rateDif
	return IRR * 100
}
