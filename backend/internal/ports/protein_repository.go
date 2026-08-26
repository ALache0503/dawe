package ports

import (
	"context"

	"github.com/ALache0503/dawe/backend/internal/domain"
)

type ProteinRepository interface {
	List(context.Context, domain.ProteinSearch) (domain.ProteinPage, error)
	FindDetailsByAccession(context.Context, string) (domain.ProteinDetails, error)
	Create(context.Context, domain.Protein) (domain.Protein, error)
	Update(context.Context, domain.Protein) (domain.Protein, error)
	Delete(context.Context, string) error
}
