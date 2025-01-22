package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/trace"
)

func AttachCorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.Request.Header.Get(trace.CorrelationIdKey)
		if correlationId == "" {
			correlationId = trace.GenerateCorrelationId()
		}

		ctx := trace.AddCorrelationIdToLogContext(c.Request.Context(), correlationId)
		ctx = context.WithValue(ctx, trace.CorrelationIdKey, correlationId)
		c.Request = c.Request.Clone(ctx)
	}
}
