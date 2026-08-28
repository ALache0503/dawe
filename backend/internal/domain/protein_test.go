package domain

import (
	"errors"
	"testing"
)

func intPointer(value int) *int {
	return &value
}

func validProtein() Protein {
	return Protein{
		Accession:       "P69905",
		TaxonID:         9606,
		EntryName:       "HBA_HUMAN",
		ProteinName:     "Hemoglobin subunit alpha",
		Reviewed:        true,
		AnnotationScore: 5,
		Mass:            intPointer(15258),
		Length:          142,
		Sequence:        "MVLSPADKTNVKAAWGKVGAHAGEYGAEALERMF",
	}
}

func TestProteinValidate(t *testing.T) {
	tests := []struct {
		name    string
		protein Protein
		wantErr error
	}{
		{
			name:    "valid protein with mass",
			protein: validProtein(),
			wantErr: nil,
		},
		{
			name: "valid protein without optional mass",
			protein: func() Protein {
				protein := validProtein()
				protein.Mass = nil
				return protein
			}(),
			wantErr: nil,
		},
		{
			name: "missing accession",
			protein: func() Protein {
				protein := validProtein()
				protein.Accession = ""
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "accession longer than ten characters",
			protein: func() Protein {
				protein := validProtein()
				protein.Accession = "A0A0C5B5G61"
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "invalid taxon ID",
			protein: func() Protein {
				protein := validProtein()
				protein.TaxonID = 0
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "empty entry name",
			protein: func() Protein {
				protein := validProtein()
				protein.EntryName = ""
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "empty protein name",
			protein: func() Protein {
				protein := validProtein()
				protein.ProteinName = ""
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "annotation score below minimum",
			protein: func() Protein {
				protein := validProtein()
				protein.AnnotationScore = 0
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "annotation score above maximum",
			protein: func() Protein {
				protein := validProtein()
				protein.AnnotationScore = 6
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "zero length",
			protein: func() Protein {
				protein := validProtein()
				protein.Length = 0
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "empty sequence",
			protein: func() Protein {
				protein := validProtein()
				protein.Sequence = ""
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "zero mass is invalid when mass is provided",
			protein: func() Protein {
				protein := validProtein()
				protein.Mass = intPointer(0)
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
		{
			name: "negative mass is invalid when mass is provided",
			protein: func() Protein {
				protein := validProtein()
				protein.Mass = intPointer(-100)
				return protein
			}(),
			wantErr: ErrInvalidInput,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.protein.Validate()

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}
