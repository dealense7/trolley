package product

import (
	"fmt"
	"net/http"
	"storePrices/internal/platform/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	repo *Repository
	log  *zap.Logger
}

func NewHandler(repo *Repository, logFactory logger.Factory) *Handler {
	return &Handler{
		repo: repo,
		log:  logFactory.For(logger.Product),
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/products")
	g.GET("/:id", h.GetProduct)
}

func (h *Handler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	prod, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		logMsg := fmt.Sprintf("failed to fetch %d: %s", id, err.Error())
		h.log.Error(logMsg, zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, prod)
}
