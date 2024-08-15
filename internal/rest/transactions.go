package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	. "github.com/julienlavocat/spacetraders/.gen/spacetraders/public/table"
	"github.com/julienlavocat/spacetraders/internal/rest/adapters"
)

func listTransaction(c *gin.Context) {
	var params struct {
		CorrelationId string `form:"correlationId"`
		Page          int    `form:"page,default=1"`
		Limit         int    `form:"limit,default=20"`
		Since         int    `form:"since,default=300"`
	}
	err := c.Bind(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	q := Transactions.SELECT(Transactions.AllColumns).LIMIT(int64(params.Limit)).OFFSET(int64((params.Page - 1) * params.Limit)).ORDER_BY(Transactions.Timestamp.DESC())

	if params.CorrelationId != "" {
		q.WHERE(Transactions.CorrelationID.EQ(String(params.CorrelationId)))
	}

	var results []model.Transactions
	if err = q.Query(db, &results); err != nil {
		internalServerError(c, "unable to query transactions", err)
		return
	}

	c.JSON(http.StatusOK, adapters.AdaptTransactions(results))
}
