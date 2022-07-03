package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PaybackPeriodRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/PP")
	ping.POST("/", calculatepaybackPeriod)
}

type PaybackPeriod struct {
	Investment float64   `json:"investment"`
	Period     int       `json:"period"`
	Cashflows  []float64 `json:"cashflows"`
}

func calculatepaybackPeriod(c *gin.Context) {
	var (
		pp                  PaybackPeriod
		paybackPeriod       float64
		accumulatedCashflow float64
		cashBeforePeriod    float64
	)
	c.ShouldBindJSON(&pp)

	paybackPeriod, accumulatedCashflow = calculatePeriod(pp.Investment, pp.Cashflows)
	for i := 0; i < int(paybackPeriod); i++ {
		cashBeforePeriod += pp.Cashflows[i]
	}
	//fmt.Printf("Cash before period %.2f is  %.2f \n", paybackPeriod, cashBeforePeriod)
	//fmt.Printf("Accumulated cashflow: %.2f \n", accumulatedCashflow)
	paybackPeriod = paybackPeriod + ((pp.Investment - cashBeforePeriod) / (accumulatedCashflow - cashBeforePeriod))
	//fmt.Printf("Payback period in %.2f \n", paybackPeriod)

	c.JSON(http.StatusOK, gin.H{
		"status":              "success",
		"paybackPeriod":       paybackPeriod,
		"accumulatedCashflow": accumulatedCashflow,
		"cashBeforePeriod":    cashBeforePeriod,
	})
}

func calculatePeriod(investment float64, cashflow []float64) (period, accumulatedCashflow float64) {
	// iterate over cashflow array and calculate accumulated cashflow
	// if accumulated cashflow is greater than investment, return both current period
	// if accumulated cashflow is less than investment, return current period
	for i := 0; i < len(cashflow); i++ {
		accumulatedCashflow += cashflow[i]
		if accumulatedCashflow > investment {
			return float64(i), accumulatedCashflow
		}
	}
	return 0, 0
}
