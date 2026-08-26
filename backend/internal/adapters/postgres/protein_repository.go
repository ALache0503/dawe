package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ALache0503/dawe/backend/internal/domain"
)

type ProteinRepository struct {
	pool *pgxpool.Pool
}

func NewProteinRepository(pool *pgxpool.Pool) *ProteinRepository {
	return &ProteinRepository{pool: pool}
}

func (r *ProteinRepository) List(ctx context.Context, search domain.ProteinSearch) (domain.ProteinPage, error) {
	const whereClause = `
        WHERE (
            $1 = '' OR
            p.accession ILIKE '%' || $1 || '%' OR
            p.entry_name ILIKE '%' || $1 || '%' OR
            p.protein_name ILIKE '%' || $1 || '%'
        )
        AND ($2::boolean IS NULL OR p.reviewed = $2)
        AND ($3::integer IS NULL OR p.taxon_id = $3)`

	var total int64
	countSQL := `SELECT COUNT(*) FROM protein p ` + whereClause
	if err := r.pool.QueryRow(ctx, countSQL, search.Query, search.Reviewed, search.TaxonID).Scan(&total); err != nil {
		return domain.ProteinPage{}, fmt.Errorf("count proteins: %w", err)
	}

	offset := (search.Page - 1) * search.PageSize
	listSQL := `
        SELECT
            p.accession,
            p.entry_name,
            p.protein_name,
            o.common_name,
            p.reviewed,
            p.annotation_score,
            p.length
        FROM protein p
        JOIN organism o ON o.taxon_id = p.taxon_id
    ` + whereClause + `
        ORDER BY p.accession ASC
        LIMIT $4 OFFSET $5`

	rows, err := r.pool.Query(ctx, listSQL,
		search.Query,
		search.Reviewed,
		search.TaxonID,
		search.PageSize,
		offset,
	)
	if err != nil {
		return domain.ProteinPage{}, fmt.Errorf("list proteins: %w", err)
	}
	defer rows.Close()

	items := make([]domain.ProteinListItem, 0)
	for rows.Next() {
		var item domain.ProteinListItem
		if err := rows.Scan(
			&item.Accession,
			&item.EntryName,
			&item.ProteinName,
			&item.OrganismName,
			&item.Reviewed,
			&item.AnnotationScore,
			&item.Length,
		); err != nil {
			return domain.ProteinPage{}, fmt.Errorf("scan protein list item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.ProteinPage{}, fmt.Errorf("iterate protein list: %w", err)
	}

	return domain.ProteinPage{
		Items:      items,
		Page:       search.Page,
		PageSize:   search.PageSize,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(search.PageSize))),
	}, nil
}

func (r *ProteinRepository) FindDetailsByAccession(ctx context.Context, accession string) (domain.ProteinDetails, error) {
	details := domain.ProteinDetails{
		Comments:        make([]domain.ProteinComment, 0),
		Features:        make([]domain.ProteinFeature, 0),
		GoTerms:         make([]domain.GoTerm, 0),
		Keywords:        make([]domain.Keyword, 0),
		CrossReferences: make([]domain.CrossReference, 0),
	}

	const proteinSQL = `
        SELECT
            p.accession, p.taxon_id, p.entry_name, p.protein_name,
            p.reviewed, p.annotation_score, p.mass, p.length, p.sequence,
            p.protein_existence, p.gene_names,
            o.taxon_id, o.common_name, o.scientific_name
        FROM protein p
        JOIN organism o ON o.taxon_id = p.taxon_id
        WHERE p.accession = $1`

	err := r.pool.QueryRow(ctx, proteinSQL, accession).Scan(
		&details.Protein.Accession,
		&details.Protein.TaxonID,
		&details.Protein.EntryName,
		&details.Protein.ProteinName,
		&details.Protein.Reviewed,
		&details.Protein.AnnotationScore,
		&details.Protein.Mass,
		&details.Protein.Length,
		&details.Protein.Sequence,
		&details.Protein.ProteinExistence,
		&details.Protein.GeneNames,
		&details.Organism.TaxonID,
		&details.Organism.CommonName,
		&details.Organism.ScientificName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ProteinDetails{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.ProteinDetails{}, fmt.Errorf("find protein: %w", err)
	}

	var loadErr error
	if details.Comments, loadErr = r.loadComments(ctx, accession); loadErr != nil {
		return domain.ProteinDetails{}, loadErr
	}
	if details.Features, loadErr = r.loadFeatures(ctx, accession); loadErr != nil {
		return domain.ProteinDetails{}, loadErr
	}
	if details.GoTerms, loadErr = r.loadGoTerms(ctx, accession); loadErr != nil {
		return domain.ProteinDetails{}, loadErr
	}
	if details.Keywords, loadErr = r.loadKeywords(ctx, accession); loadErr != nil {
		return domain.ProteinDetails{}, loadErr
	}
	if details.CrossReferences, loadErr = r.loadCrossReferences(ctx, accession); loadErr != nil {
		return domain.ProteinDetails{}, loadErr
	}

	return details, nil
}

func (r *ProteinRepository) Create(ctx context.Context, protein domain.Protein) (domain.Protein, error) {
	const sql = `
        INSERT INTO protein (
            accession, taxon_id, entry_name, protein_name, reviewed,
            annotation_score, mass, length, sequence, protein_existence, gene_names
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.pool.Exec(ctx, sql,
		protein.Accession, protein.TaxonID, protein.EntryName, protein.ProteinName,
		protein.Reviewed, protein.AnnotationScore, protein.Mass, protein.Length,
		protein.Sequence, protein.ProteinExistence, protein.GeneNames,
	)
	if err != nil {
		return domain.Protein{}, mapPostgresError(err)
	}
	return protein, nil
}

func (r *ProteinRepository) Update(ctx context.Context, protein domain.Protein) (domain.Protein, error) {
	const sql = `
        UPDATE protein
        SET
            taxon_id = $2,
            entry_name = $3,
            protein_name = $4,
            reviewed = $5,
            annotation_score = $6,
            mass = $7,
            length = $8,
            sequence = $9,
            protein_existence = $10,
            gene_names = $11
        WHERE accession = $1`

	result, err := r.pool.Exec(ctx, sql,
		protein.Accession, protein.TaxonID, protein.EntryName, protein.ProteinName,
		protein.Reviewed, protein.AnnotationScore, protein.Mass, protein.Length,
		protein.Sequence, protein.ProteinExistence, protein.GeneNames,
	)
	if err != nil {
		return domain.Protein{}, mapPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return domain.Protein{}, domain.ErrNotFound
	}
	return protein, nil
}

func (r *ProteinRepository) Delete(ctx context.Context, accession string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM protein WHERE accession = $1`, accession)
	if err != nil {
		return fmt.Errorf("delete protein: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ProteinRepository) loadComments(ctx context.Context, accession string) ([]domain.ProteinComment, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT pc.comment_id, ct.type_name, COALESCE(pc.comment_text, '')
        FROM proteincomment pc
        JOIN commenttype ct ON ct.type_id = pc.commenttype_id
        WHERE pc.accession = $1
        ORDER BY pc.comment_id`, accession)
	if err != nil {
		return nil, fmt.Errorf("load comments: %w", err)
	}
	defer rows.Close()

	result := make([]domain.ProteinComment, 0)
	for rows.Next() {
		var item domain.ProteinComment
		if err := rows.Scan(&item.CommentID, &item.TypeName, &item.CommentText); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *ProteinRepository) loadFeatures(ctx context.Context, accession string) ([]domain.ProteinFeature, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT pf.feature_id, ft.type_name, COALESCE(pf.description, ''),
               pf.start_position, pf.end_position
        FROM proteinfeature pf
        JOIN featuretype ft ON ft.type_id = pf.featuretype_id
        WHERE pf.accession = $1
        ORDER BY pf.start_position, pf.feature_id`, accession)
	if err != nil {
		return nil, fmt.Errorf("load features: %w", err)
	}
	defer rows.Close()

	result := make([]domain.ProteinFeature, 0)
	for rows.Next() {
		var item domain.ProteinFeature
		if err := rows.Scan(&item.FeatureID, &item.TypeName, &item.Description, &item.StartPosition, &item.EndPosition); err != nil {
			return nil, fmt.Errorf("scan feature: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *ProteinRepository) loadGoTerms(ctx context.Context, accession string) ([]domain.GoTerm, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT gt.go_id, COALESCE(gt.term_name, ''), COALESCE(gt.category, ''), gt.aspect
        FROM annotated_with aw
        JOIN goterm gt ON gt.go_id = aw.go_id
        WHERE aw.accession = $1
        ORDER BY gt.go_id`, accession)
	if err != nil {
		return nil, fmt.Errorf("load GO terms: %w", err)
	}
	defer rows.Close()

	result := make([]domain.GoTerm, 0)
	for rows.Next() {
		var item domain.GoTerm
		if err := rows.Scan(&item.GoID, &item.TermName, &item.Category, &item.Aspect); err != nil {
			return nil, fmt.Errorf("scan GO term: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *ProteinRepository) loadKeywords(ctx context.Context, accession string) ([]domain.Keyword, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT k.keyword_id, k.keyword_name
        FROM tagged_with tw
        JOIN keyword k ON k.keyword_id = tw.keyword_id
        WHERE tw.accession = $1
        ORDER BY k.keyword_name`, accession)
	if err != nil {
		return nil, fmt.Errorf("load keywords: %w", err)
	}
	defer rows.Close()

	result := make([]domain.Keyword, 0)
	for rows.Next() {
		var item domain.Keyword
		if err := rows.Scan(&item.KeywordID, &item.KeywordName); err != nil {
			return nil, fmt.Errorf("scan keyword: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *ProteinRepository) loadCrossReferences(ctx context.Context, accession string) ([]domain.CrossReference, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT reference_id, database_name
        FROM crossreference
        WHERE accession = $1
        ORDER BY database_name, reference_id`, accession)
	if err != nil {
		return nil, fmt.Errorf("load cross references: %w", err)
	}
	defer rows.Close()

	result := make([]domain.CrossReference, 0)
	for rows.Next() {
		var item domain.CrossReference
		if err := rows.Scan(&item.ReferenceID, &item.DatabaseName); err != nil {
			return nil, fmt.Errorf("scan cross reference: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func mapPostgresError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return domain.ErrConflict
		case "23503", "23514", "22001", "22P02":
			return domain.ErrInvalidInput
		}
	}
	return fmt.Errorf("database operation: %w", err)
}
