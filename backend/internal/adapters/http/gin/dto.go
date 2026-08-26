package ginadapter

import "github.com/ALache0503/dawe/backend/internal/domain"

type createProteinRequest struct {
	Accession        string  `json:"accession" binding:"required,max=10"`
	TaxonID          int     `json:"taxonId" binding:"required,gt=0"`
	EntryName        string  `json:"entryName" binding:"required,max=50"`
	ProteinName      string  `json:"proteinName" binding:"required"`
	Reviewed         bool    `json:"reviewed"`
	AnnotationScore  int16   `json:"annotationScore" binding:"required,min=1,max=5"`
	Mass             *int    `json:"mass"`
	Length           int     `json:"length" binding:"required,gt=0"`
	Sequence         string  `json:"sequence" binding:"required"`
	ProteinExistence *string `json:"proteinExistence"`
	GeneNames        *string `json:"geneNames"`
}

type updateProteinRequest struct {
	TaxonID          int     `json:"taxonId" binding:"required,gt=0"`
	EntryName        string  `json:"entryName" binding:"required,max=50"`
	ProteinName      string  `json:"proteinName" binding:"required"`
	Reviewed         bool    `json:"reviewed"`
	AnnotationScore  int16   `json:"annotationScore" binding:"required,min=1,max=5"`
	Mass             *int    `json:"mass"`
	Length           int     `json:"length" binding:"required,gt=0"`
	Sequence         string  `json:"sequence" binding:"required"`
	ProteinExistence *string `json:"proteinExistence"`
	GeneNames        *string `json:"geneNames"`
}

func (r createProteinRequest) toDomain() domain.Protein {
	return domain.Protein{
		Accession: r.Accession, TaxonID: r.TaxonID, EntryName: r.EntryName,
		ProteinName: r.ProteinName, Reviewed: r.Reviewed,
		AnnotationScore: r.AnnotationScore, Mass: r.Mass, Length: r.Length,
		Sequence: r.Sequence, ProteinExistence: r.ProteinExistence, GeneNames: r.GeneNames,
	}
}

func (r updateProteinRequest) toDomain(accession string) domain.Protein {
	return domain.Protein{
		Accession: accession, TaxonID: r.TaxonID, EntryName: r.EntryName,
		ProteinName: r.ProteinName, Reviewed: r.Reviewed,
		AnnotationScore: r.AnnotationScore, Mass: r.Mass, Length: r.Length,
		Sequence: r.Sequence, ProteinExistence: r.ProteinExistence, GeneNames: r.GeneNames,
	}
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
