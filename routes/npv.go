package routes

import (
	"fmt"
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
	ping.POST("/", getNPV)
}

func getNPV(c *gin.Context) {
	var npv NPV
	var flag bool
	c.ShouldBindJSON(&npv)
	result, _ := calculateNPV(npv.Investment, npv.Cashflows, npv.InterestRate/100)
	fmt.Printf("NPV is: %.2f \n", result)
	if result > 0 {
		flag = true
	} else {
		flag = false
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"NPV":    result,
		"flag":   flag,
	})
}

func calculateNPV(investment int, cashflows []int, interestRate float64) (float64, []float64) {
	// calculate NPV
	var presentValue []float64
	var result float64

	for i := 0; i < len(cashflows); i++ {
		interest := 1 / math.Pow(1+interestRate, float64(i+1))
		fmt.Println(interest)
		result = float64(cashflows[i]) * interest
		presentValue = append(presentValue, result)
	}
	// print present value
	fmt.Printf("Present value is: %.2f \n", presentValue)
	// sum all values in presentValue array
	var sum float64
	for i := 0; i < len(presentValue); i++ {
		sum += presentValue[i]
	}
	fmt.Printf("Sum of present value is: %.2f \n\n", sum)
	return sum - float64(investment), presentValue
}
