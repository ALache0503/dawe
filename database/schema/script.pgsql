DROP TABLE IF EXISTS ProteinFeature, ProteinComment, CrossReference,
    annotated_with, tagged_with, Protein, GoTerm, Keyword,
    CommentType, FeatureType, Organism CASCADE;


CREATE TABLE Keyword (
    keyword_id VARCHAR(10) NOT NULL PRIMARY KEY,
    keyword_name VARCHAR(255) NOT NULL
);

CREATE TABLE Organism (
    taxon_id INTEGER NOT NULL PRIMARY KEY,
    common_name VARCHAR(100) NOT NULL,
    scientific_name VARCHAR(100) NOT NULL
);

CREATE TABLE CommentType (
    type_id INTEGER GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL
);

CREATE TABLE GoTerm (
    go_id VARCHAR(10) NOT NULL PRIMARY KEY,
    term_name TEXT,
    category VARCHAR(50),
    aspect VARCHAR(1) NOT NULL CHECK (aspect IN ('P', 'F', 'C'))
);

CREATE TABLE FeatureType (
    type_id INTEGER GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL
);

CREATE TABLE Protein (
    accession VARCHAR(10) NOT NULL PRIMARY KEY,
    taxon_id INTEGER REFERENCES Organism(taxon_id) NOT NULL,
    entry_name VARCHAR(50) NOT NULL,
    protein_name TEXT NOT NULL,
    reviewed BOOLEAN NOT NULL,
    annotation_score SMALLINT NOT NULL CHECK (annotation_score BETWEEN 1 AND 5),
    mass INTEGER,
    length INTEGER NOT NULL,
    sequence TEXT NOT NULL,
    protein_existence VARCHAR(50),
    gene_names TEXT
);

CREATE TABLE ProteinComment (
    accession VARCHAR(10) NOT NULL REFERENCES Protein(accession) ON DELETE CASCADE,
    comment_id INTEGER NOT NULL,
    commentType_id INTEGER REFERENCES CommentType(type_id) NOT NULL,
    comment_text TEXT,
    PRIMARY KEY (accession, comment_id)
);

CREATE TABLE CrossReference (
    accession VARCHAR(10) NOT NULL REFERENCES Protein(accession) ON DELETE CASCADE,
    reference_id VARCHAR(100) NOT NULL,
    database_name VARCHAR(50) NOT NULL,
    PRIMARY KEY (accession, reference_id)
);

CREATE TABLE annotated_with (
    accession VARCHAR(10) NOT NULL REFERENCES Protein(accession) ON DELETE CASCADE,
    go_id VARCHAR(10) NOT NULL REFERENCES GoTerm(go_id) ON DELETE CASCADE,
    PRIMARY KEY (accession, go_id)
);

CREATE TABLE tagged_with (
    accession VARCHAR(10) NOT NULL REFERENCES Protein(accession) ON DELETE CASCADE,
    keyword_id VARCHAR(10) NOT NULL REFERENCES Keyword(keyword_id) ON DELETE CASCADE,
    PRIMARY KEY (accession, keyword_id)
);

CREATE TABLE ProteinFeature (
    accession VARCHAR(10) NOT NULL REFERENCES Protein(accession) ON DELETE CASCADE,
    feature_id INTEGER NOT NULL,
    featureType_id INTEGER NOT NULL REFERENCES FeatureType(type_id),
    description TEXT,
    start_position INTEGER NOT NULL,
    end_position INTEGER NOT NULL,
    PRIMARY KEY (accession, feature_id)
);
