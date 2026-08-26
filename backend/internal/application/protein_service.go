package application

import (
	"context"
	"strings"

	"github.com/ALache0503/dawe/backend/internal/domain"
	"github.com/ALache0503/dawe/backend/internal/ports"
)

type ProteinService struct {
	repository ports.ProteinRepository
}

func NewProteinService(repository ports.ProteinRepository) *ProteinService {
	return &ProteinService{repository: repository}
}

func (s *ProteinService) List(ctx context.Context, search domain.ProteinSearch) (domain.ProteinPage, error) {
	if search.Page < 1 {
		return domain.ProteinPage{}, domain.ErrInvalidPage
	}
	if search.PageSize < 1 || search.PageSize > 100 {
		return domain.ProteinPage{}, domain.ErrInvalidPageSize
	}

	search.Query = strings.TrimSpace(search.Query)
	return s.repository.List(ctx, search)
}

func (s *ProteinService) GetDetails(ctx context.Context, accession string) (domain.ProteinDetails, error) {
	accession = strings.TrimSpace(accession)
	if accession == "" {
		return domain.ProteinDetails{}, domain.ErrInvalidInput
	}
	return s.repository.FindDetailsByAccession(ctx, accession)
}

func (s *ProteinService) Create(ctx context.Context, protein domain.Protein) (domain.Protein, error) {
	protein.Accession = strings.TrimSpace(protein.Accession)
	protein.EntryName = strings.TrimSpace(protein.EntryName)
	protein.ProteinName = strings.TrimSpace(protein.ProteinName)
	protein.Sequence = strings.TrimSpace(protein.Sequence)

	if err := protein.Validate(); err != nil {
		return domain.Protein{}, err
	}
	return s.repository.Create(ctx, protein)
}

func (s *ProteinService) Update(ctx context.Context, accession string, protein domain.Protein) (domain.Protein, error) {
	accession = strings.TrimSpace(accession)
	if accession == "" {
		return domain.Protein{}, domain.ErrInvalidInput
	}

	protein.Accession = accession
	protein.EntryName = strings.TrimSpace(protein.EntryName)
	protein.ProteinName = strings.TrimSpace(protein.ProteinName)
	protein.Sequence = strings.TrimSpace(protein.Sequence)

	if err := protein.Validate(); err != nil {
		return domain.Protein{}, err
	}
	return s.repository.Update(ctx, protein)
}

func (s *ProteinService) Delete(ctx context.Context, accession string) error {
	accession = strings.TrimSpace(accession)
	if accession == "" {
		return domain.ErrInvalidInput
	}
	return s.repository.Delete(ctx, accession)
}
