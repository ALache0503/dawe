package ginadapter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ALache0503/dawe/backend/internal/application"
	"github.com/ALache0503/dawe/backend/internal/domain"
)

type handlerFakeRepository struct {
	lastSearch    domain.ProteinSearch
	lastAccession string
	listResult    domain.ProteinPage
	detailsResult domain.ProteinDetails
	listError     error
	detailsError  error
	createError   error
	updateError   error
	deleteError   error
}

func (f *handlerFakeRepository) List(
	_ context.Context,
	search domain.ProteinSearch,
) (domain.ProteinPage, error) {
	f.lastSearch = search
	return f.listResult, f.listError
}

func (f *handlerFakeRepository) FindDetailsByAccession(
	_ context.Context,
	accession string,
) (domain.ProteinDetails, error) {
	f.lastAccession = accession
	return f.detailsResult, f.detailsError
}

func (f *handlerFakeRepository) Create(
	_ context.Context,
	protein domain.Protein,
) (domain.Protein, error) {
	if f.createError != nil {
		return domain.Protein{}, f.createError
	}

	return protein, nil
}

func (f *handlerFakeRepository) Update(
	_ context.Context,
	protein domain.Protein,
) (domain.Protein, error) {
	if f.updateError != nil {
		return domain.Protein{}, f.updateError
	}

	return protein, nil
}

func (f *handlerFakeRepository) Delete(
	_ context.Context,
	accession string,
) error {
	f.lastAccession = accession
	return f.deleteError
}

func newTestRouter(repository *handlerFakeRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	service := application.NewProteinService(repository)
	handler := NewProteinHandler(service)

	router := gin.New()

	router.GET("/proteins", handler.List)
	router.GET("/proteins/:accession", handler.GetDetails)
	router.POST("/proteins", handler.Create)
	router.PUT("/proteins/:accession", handler.Update)
	router.DELETE("/proteins/:accession", handler.Delete)

	return router
}

func TestProteinHandlerListReturnsPage(t *testing.T) {
	repository := &handlerFakeRepository{
		listResult: domain.ProteinPage{
			Items: []domain.ProteinListItem{
				{
					Accession:       "P69905",
					EntryName:       "HBA_HUMAN",
					ProteinName:     "Hemoglobin subunit alpha",
					OrganismName:    "Human",
					Reviewed:        true,
					AnnotationScore: 5,
					Length:          142,
				},
			},
			Page:       2,
			PageSize:   10,
			TotalItems: 11,
			TotalPages: 2,
		},
	}

	router := newTestRouter(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/proteins?page=2&pageSize=10&search=hemoglobin&reviewed=true&taxonId=9606",
		nil,
	)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", response.Code, response.Body.String())
	}

	if repository.lastSearch.Page != 2 {
		t.Fatalf("expected page 2, got %d", repository.lastSearch.Page)
	}

	if repository.lastSearch.PageSize != 10 {
		t.Fatalf("expected page size 10, got %d", repository.lastSearch.PageSize)
	}

	if repository.lastSearch.Query != "hemoglobin" {
		t.Fatalf("expected search hemoglobin, got %q", repository.lastSearch.Query)
	}

	if repository.lastSearch.Reviewed == nil || !*repository.lastSearch.Reviewed {
		t.Fatal("expected reviewed filter true")
	}

	if repository.lastSearch.TaxonID == nil || *repository.lastSearch.TaxonID != 9606 {
		t.Fatal("expected taxon ID filter 9606")
	}

	var page domain.ProteinPage
	if err := json.Unmarshal(response.Body.Bytes(), &page); err != nil {
		t.Fatalf("could not decode JSON response: %v", err)
	}

	if page.TotalItems != 11 {
		t.Fatalf("expected totalItems 11, got %d", page.TotalItems)
	}

	if len(page.Items) != 1 {
		t.Fatalf("expected 1 protein item, got %d", len(page.Items))
	}
}

func TestProteinHandlerListUsesDefaults(t *testing.T) {
	repository := &handlerFakeRepository{
		listResult: domain.ProteinPage{
			Items:    []domain.ProteinListItem{},
			Page:     1,
			PageSize: 20,
		},
	}

	router := newTestRouter(repository)
	request := httptest.NewRequest(http.MethodGet, "/proteins", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	if repository.lastSearch.Page != 1 {
		t.Fatalf("expected default page 1, got %d", repository.lastSearch.Page)
	}

	if repository.lastSearch.PageSize != 20 {
		t.Fatalf("expected default page size 20, got %d", repository.lastSearch.PageSize)
	}
}

func TestProteinHandlerListRejectsInvalidQueryParameters(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "page is zero",
			url:  "/proteins?page=0",
		},
		{
			name: "page is text",
			url:  "/proteins?page=abc",
		},
		{
			name: "page size is too large",
			url:  "/proteins?pageSize=101",
		},
		{
			name: "reviewed is invalid",
			url:  "/proteins?reviewed=yes",
		},
		{
			name: "taxon ID is invalid",
			url:  "/proteins?taxonId=abc",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &handlerFakeRepository{}
			router := newTestRouter(repository)

			request := httptest.NewRequest(http.MethodGet, test.url, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf(
					"expected status 400, got %d; body: %s",
					response.Code,
					response.Body.String(),
				)
			}
		})
	}
}

func TestProteinHandlerGetDetailsReturnsNotFound(t *testing.T) {
	repository := &handlerFakeRepository{
		detailsError: domain.ErrNotFound,
	}

	router := newTestRouter(repository)
	request := httptest.NewRequest(http.MethodGet, "/proteins/UNKNOWN", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d; body: %s",
			response.Code,
			response.Body.String(),
		)
	}

	var apiError errorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &apiError); err != nil {
		t.Fatalf("could not decode JSON error response: %v", err)
	}

	if apiError.Code != "not_found" {
		t.Fatalf("expected error code not_found, got %q", apiError.Code)
	}
}

func TestProteinHandlerGetDetailsReturnsProtein(t *testing.T) {
	repository := &handlerFakeRepository{
		detailsResult: domain.ProteinDetails{
			Protein: domain.Protein{
				Accession:   "P69905",
				ProteinName: "Hemoglobin subunit alpha",
			},
			Organism: domain.Organism{
				TaxonID:        9606,
				CommonName:     "Human",
				ScientificName: "Homo sapiens",
			},
			Comments:        []domain.ProteinComment{},
			Features:        []domain.ProteinFeature{},
			GoTerms:         []domain.GoTerm{},
			Keywords:        []domain.Keyword{},
			CrossReferences: []domain.CrossReference{},
		},
	}

	router := newTestRouter(repository)
	request := httptest.NewRequest(http.MethodGet, "/proteins/P69905", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d; body: %s",
			response.Code,
			response.Body.String(),
		)
	}

	if repository.lastAccession != "P69905" {
		t.Fatalf("expected accession P69905, got %q", repository.lastAccession)
	}

	var details domain.ProteinDetails
	if err := json.Unmarshal(response.Body.Bytes(), &details); err != nil {
		t.Fatalf("could not decode JSON response: %v", err)
	}

	if details.Protein.Accession != "P69905" {
		t.Fatalf(
			"expected protein P69905, got %q",
			details.Protein.Accession,
		)
	}
}

func TestProteinHandlerCreateReturnsCreated(t *testing.T) {
	repository := &handlerFakeRepository{}
	router := newTestRouter(repository)

	body := `{
		"accession": "P69905",
		"taxonId": 9606,
		"entryName": "HBA_HUMAN",
		"proteinName": "Hemoglobin subunit alpha",
		"reviewed": true,
		"annotationScore": 5,
		"mass": 15258,
		"length": 142,
		"sequence": "MVLSPADKTNVKAAWGKVGAHAGEYGAEALERMF",
		"proteinExistence": "Evidence at protein level",
		"geneNames": "HBA1 HBA2"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/proteins",
		strings.NewReader(body),
	)
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"expected status 201, got %d; body: %s",
			response.Code,
			response.Body.String(),
		)
	}

	if location := response.Header().Get("Location"); location != "/api/v1/proteins/P69905" {
		t.Fatalf("unexpected Location header: %q", location)
	}
}

func TestProteinHandlerCreateRejectsInvalidJSON(t *testing.T) {
	repository := &handlerFakeRepository{}
	router := newTestRouter(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/proteins",
		strings.NewReader(`{"accession":`),
	)
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestProteinHandlerDeleteReturnsNoContent(t *testing.T) {
	repository := &handlerFakeRepository{}
	router := newTestRouter(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/proteins/P69905",
		nil,
	)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status 204, got %d; body: %s",
			response.Code,
			response.Body.String(),
		)
	}

	if repository.lastAccession != "P69905" {
		t.Fatalf("expected accession P69905, got %q", repository.lastAccession)
	}
}

func TestProteinHandlerDeleteReturnsNotFound(t *testing.T) {
	repository := &handlerFakeRepository{
		deleteError: domain.ErrNotFound,
	}
	router := newTestRouter(repository)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/proteins/UNKNOWN",
		nil,
	)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d; body: %s",
			response.Code,
			response.Body.String(),
		)
	}
}

func TestProteinHandlerMapsUnexpectedErrorToInternalServerError(t *testing.T) {
	repository := &handlerFakeRepository{
		listError: errors.New("unexpected database error"),
	}
	router := newTestRouter(repository)

	request := httptest.NewRequest(http.MethodGet, "/proteins", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d; body: %s",
			response.Code,
			response.Body.String(),
		)
	}
}
