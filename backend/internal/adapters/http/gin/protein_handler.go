package ginadapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ALache0503/dawe/backend/internal/application"
	"github.com/ALache0503/dawe/backend/internal/domain"
)

type ProteinHandler struct {
	service *application.ProteinService
}

func NewProteinHandler(service *application.ProteinService) *ProteinHandler {
	return &ProteinHandler{service: service}
}

func (h *ProteinHandler) List(c *gin.Context) {
	page, err := optionalPositiveInt(c.Query("page"), 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{"validation_error", "page must be a positive integer."})
		return
	}
	pageSize, err := optionalPositiveInt(c.Query("pageSize"), 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{"validation_error", "pageSize must be a positive integer."})
		return
	}

	search := domain.ProteinSearch{
		Query: c.Query("search"), Page: page, PageSize: pageSize,
	}
	if reviewed, exists, err := optionalBool(c.Query("reviewed")); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{"validation_error", "reviewed must be true or false."})
		return
	} else if exists {
		search.Reviewed = &reviewed
	}
	if taxonID, exists, err := optionalInt(c.Query("taxonId")); err != nil || (exists && taxonID <= 0) {
		c.JSON(http.StatusBadRequest, errorResponse{"validation_error", "taxonId must be a positive integer."})
		return
	} else if exists {
		search.TaxonID = &taxonID
	}

	result, err := h.service.List(c.Request.Context(), search)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ProteinHandler) GetDetails(c *gin.Context) {
	result, err := h.service.GetDetails(c.Request.Context(), c.Param("accession"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ProteinHandler) Create(c *gin.Context) {
	var request createProteinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{"validation_error", "The request body is invalid."})
		return
	}

	protein, err := h.service.Create(c.Request.Context(), request.toDomain())
	if err != nil {
		writeError(c, err)
		return
	}
	c.Header("Location", "/api/v1/proteins/"+protein.Accession)
	c.JSON(http.StatusCreated, protein)
}

func (h *ProteinHandler) Update(c *gin.Context) {
	var request updateProteinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{"validation_error", "The request body is invalid."})
		return
	}

	protein, err := h.service.Update(c.Request.Context(), c.Param("accession"), request.toDomain(c.Param("accession")))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, protein)
}

func (h *ProteinHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("accession")); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func optionalPositiveInt(raw string, defaultValue int) (int, error) {
	if raw == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, domain.ErrInvalidInput
	}
	return value, nil
}

func optionalInt(raw string) (int, bool, error) {
	if raw == "" {
		return 0, false, nil
	}
	value, err := strconv.Atoi(raw)
	return value, true, err
}

func optionalBool(raw string) (bool, bool, error) {
	if raw == "" {
		return false, false, nil
	}
	value, err := strconv.ParseBool(raw)
	return value, true, err
}
