package ginadapter

import (
	"errors"
	"net/http"

	"github.com/ALache0503/dawe/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResponse{"not_found", "Resource was not found."})
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, errorResponse{"conflict", "A record with this identifier already exists."})
	case errors.Is(err, domain.ErrInvalidInput),
		errors.Is(err, domain.ErrInvalidPage),
		errors.Is(err, domain.ErrInvalidPageSize):
		c.JSON(http.StatusBadRequest, errorResponse{"validation_error", "The request is invalid."})
	default:
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errorResponse{"internal_error", "An internal error occurred."})
	}
}
