package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ALache0503/dawe/backend/internal/application"
	"github.com/ALache0503/dawe/backend/internal/domain"
)

type fakeProteinRepository struct {
	receivedSearch domain.ProteinSearch
}

func (f *fakeProteinRepository) List(_ context.Context, search domain.ProteinSearch) (domain.ProteinPage, error) {
	f.receivedSearch = search
	return domain.ProteinPage{Page: search.Page, PageSize: search.PageSize}, nil
}

func (f *fakeProteinRepository) FindDetailsByAccession(_ context.Context, _ string) (domain.ProteinDetails, error) {
	return domain.ProteinDetails{}, nil
}

func (f *fakeProteinRepository) Create(_ context.Context, protein domain.Protein) (domain.Protein, error) {
	return protein, nil
}

func (f *fakeProteinRepository) Update(_ context.Context, protein domain.Protein) (domain.Protein, error) {
	return protein, nil
}

func (f *fakeProteinRepository) Delete(_ context.Context, _ string) error {
	return nil
}

func TestListRejectsPageSizeAboveLimit(t *testing.T) {
	repository := &fakeProteinRepository{}
	service := application.NewProteinService(repository)

	_, err := service.List(context.Background(), domain.ProteinSearch{Page: 1, PageSize: 101})
	if !errors.Is(err, domain.ErrInvalidPageSize) {
		t.Fatalf("expected ErrInvalidPageSize, got %v", err)
	}
}

func TestListDelegatesValidQuery(t *testing.T) {
	repository := &fakeProteinRepository{}
	service := application.NewProteinService(repository)

	_, err := service.List(context.Background(), domain.ProteinSearch{
		Query: " hemoglobin ", Page: 1, PageSize: 20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repository.receivedSearch.Query != "hemoglobin" {
		t.Fatalf("expected trimmed query, got %q", repository.receivedSearch.Query)
	}
}
