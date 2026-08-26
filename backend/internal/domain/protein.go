package domain

import "strings"

type Protein struct {
	Accession        string
	TaxonID          int
	EntryName        string
	ProteinName      string
	Reviewed         bool
	AnnotationScore  int16
	Mass             *int
	Length           int
	Sequence         string
	ProteinExistence *string
	GeneNames        *string
}

type Organism struct {
	TaxonID        int    `json:"taxonId"`
	CommonName     string `json:"commonName"`
	ScientificName string `json:"scientificName"`
}

type ProteinListItem struct {
	Accession       string `json:"accession"`
	EntryName       string `json:"entryName"`
	ProteinName     string `json:"proteinName"`
	OrganismName    string `json:"organismName"`
	Reviewed        bool   `json:"reviewed"`
	AnnotationScore int16  `json:"annotationScore"`
	Length          int    `json:"length"`
}

type ProteinSearch struct {
	Query    string
	Reviewed *bool
	TaxonID  *int
	Page     int
	PageSize int
}

type ProteinPage struct {
	Items      []ProteinListItem `json:"items"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalItems int64             `json:"totalItems"`
	TotalPages int               `json:"totalPages"`
}

type ProteinComment struct {
	CommentID   int    `json:"commentId"`
	TypeName    string `json:"typeName"`
	CommentText string `json:"commentText"`
}

type ProteinFeature struct {
	FeatureID     int    `json:"featureId"`
	TypeName      string `json:"typeName"`
	Description   string `json:"description"`
	StartPosition int    `json:"startPosition"`
	EndPosition   int    `json:"endPosition"`
}

type GoTerm struct {
	GoID     string `json:"goId"`
	TermName string `json:"termName"`
	Category string `json:"category"`
	Aspect   string `json:"aspect"`
}

type Keyword struct {
	KeywordID   string `json:"keywordId"`
	KeywordName string `json:"keywordName"`
}

type CrossReference struct {
	ReferenceID  string `json:"referenceId"`
	DatabaseName string `json:"databaseName"`
}

type ProteinDetails struct {
	Protein         Protein          `json:"protein"`
	Organism        Organism         `json:"organism"`
	Comments        []ProteinComment `json:"comments"`
	Features        []ProteinFeature `json:"features"`
	GoTerms         []GoTerm         `json:"goTerms"`
	Keywords        []Keyword        `json:"keywords"`
	CrossReferences []CrossReference `json:"crossReferences"`
}

func (p Protein) Validate() error {
	if strings.TrimSpace(p.Accession) == "" || len(p.Accession) > 10 {
		return ErrInvalidInput
	}
	if p.TaxonID <= 0 || strings.TrimSpace(p.EntryName) == "" || len(p.EntryName) > 50 {
		return ErrInvalidInput
	}
	if strings.TrimSpace(p.ProteinName) == "" || (p.Mass != nil && *p.Mass <= 0) || p.Length <= 0 {
		return ErrInvalidInput
	}
	if p.AnnotationScore < 1 || p.AnnotationScore > 5 || strings.TrimSpace(p.Sequence) == "" {
		return ErrInvalidInput
	}
	return nil
}
