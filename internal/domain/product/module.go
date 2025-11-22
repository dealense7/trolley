package product

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("product",
	fx.Provide(
		NewRepository,
		NewHandler,
	),
	fx.Invoke(
		func(h *Handler, r *gin.Engine) {
			h.RegisterRoutes(r)
		},
	),
)
