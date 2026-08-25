#!/usr/bin/env python3
"""
UniProtKB Data Import Script
Fetches protein data from the UniProt REST API and inserts it into a PostgreSQL database.

Prerequisites:
    pip install requests psycopg2-binary

Usage:
    python import_uniprot.py

Configuration:
    Set DB_URL and UniProt query parameters below.
"""

import requests
import psycopg2
from psycopg2.extras import execute_batch
import time
import sys

# ============================================================
# Configuration
# ============================================================

# PostgreSQL connection string — adjust to your setup
DB_URL = "postgresql://postgres:postgres@localhost:5432/uniprot_db"

# UniProt API base URL
UNIPROT_API = "https://rest.uniprot.org/uniprotkb"

# Search query: reviewed human proteins (Swiss-Prot)
# Adjust organism_id or query as needed
QUERY = "organism_id:9606 AND reviewed:true"

# Number of proteins to import (assignment says "few representative records")
PROTEIN_COUNT = 50

# Batch size for API requests (max 500 per request)
API_BATCH_SIZE = 25

# ============================================================
# Helpers
# ============================================================

ASPECT_MAP = {"P": "Biological Process", "F": "Molecular Function", "C": "Cellular Component"}


def fetch_protein(accession):
    """Fetch a single protein entry as full JSON from the UniProt API."""
    url = f"{UNIPROT_API}/{accession}.json"
    resp = requests.get(url, headers={"User-Agent": "UniProt-Importer/1.0"})
    resp.raise_for_status()
    return resp.json()


def search_proteins(query, size):
    """Search UniProt and return a list of accessions."""
    url = f"{UNIPROT_API}/search"
    params = {"query": query, "format": "json", "size": size}
    resp = requests.get(url, params=params, headers={"User-Agent": "UniProt-Importer/1.0"})
    resp.raise_for_status()
    data = resp.json()
    return [entry["primaryAccession"] for entry in data.get("results", [])]


def get_protein_name(entry):
    """Extract the recommended protein name from JSON."""
    desc = entry.get("proteinDescription", {})
    rec = desc.get("recommendedName", {})
    if rec:
        full = rec.get("fullName", {})
        return full.get("value", "Unknown")
    alt_names = desc.get("alternativeNames", [])
    if alt_names:
        return alt_names[0].get("fullName", {}).get("value", "Unknown")
    return "Unknown"


def get_gene_names(entry):
    """Extract gene names as a semicolon-separated string."""
    genes = entry.get("genes", [])
    names = []
    for g in genes:
        gn = g.get("geneName", {})
        if gn:
            names.append(gn.get("value", ""))
        for syn in g.get("synonyms", []):
            names.append(syn.get("value", ""))
    return "; ".join(names) if names else None


def get_comment_text(comment):
    """Extract text from a comment object (handles multiple formats)."""
    texts = comment.get("texts", [])
    if texts:
        return " ".join(t.get("value", "") for t in texts)
    # Some comments store text differently
    return comment.get("value", "")


def get_feature_location(feature):
    """Extract start and end positions from a feature."""
    loc = feature.get("location", {})
    start = loc.get("start", {}).get("value")
    end = loc.get("end", {}).get("value")
    return start, end


# ============================================================
# Database insert functions
# ============================================================

def insert_organism(cur, entry):
    """Insert organism if not already present."""
    org = entry.get("organism", {})
    taxon_id = org.get("taxonId")
    if not taxon_id:
        return None
    cur.execute(
        "INSERT INTO Organism (taxon_id, common_name, scientific_name) "
        "VALUES (%s, %s, %s) ON CONFLICT (taxon_id) DO NOTHING",
        (taxon_id, org.get("commonName", ""), org.get("scientificName", ""))
    )
    return taxon_id


def insert_protein(cur, entry):
    """Insert the main protein record."""
    accession = entry.get("primaryAccession")
    taxon_id = entry.get("organism", {}).get("taxonId")
    entry_name = entry.get("uniProtkbId", "")
    protein_name = get_protein_name(entry)
    reviewed = "reviewed" in entry.get("entryType", "").lower()
    annotation_score = int(entry.get("annotationScore", 1))
    seq_data = entry.get("sequence", {})
    mass = int(seq_data.get("mass", 0)) if seq_data.get("mass") else None
    length = int(seq_data.get("length", 0)) if seq_data.get("length") else None
    sequence = seq_data.get("value", "")
    protein_existence = entry.get("proteinExistence")
    gene_names = get_gene_names(entry)

    cur.execute(
        """INSERT INTO Protein
           (accession, taxon_id, entry_name, protein_name, reviewed,
            annotation_score, mass, length, sequence, protein_existence, gene_names)
           VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
           ON CONFLICT (accession) DO NOTHING""",
        (accession, taxon_id, entry_name, protein_name, reviewed,
         annotation_score, mass, length, sequence, protein_existence, gene_names)
    )
    return accession


def insert_comments(cur, entry, accession, comment_type_cache):
    """Insert all comments for a protein."""
    comments = entry.get("comments", [])
    comment_id = 1
    rows = []
    for comment in comments:
        ctype = comment.get("commentType", "UNKNOWN")
        text = get_comment_text(comment)
        if not text:
            continue

        # Resolve or create CommentType
        if ctype not in comment_type_cache:
            cur.execute(
                "INSERT INTO CommentType (type_name) VALUES (%s) RETURNING type_id",
                (ctype,)
            )
            comment_type_cache[ctype] = cur.fetchone()[0]

        type_id = comment_type_cache[ctype]
        rows.append((accession, comment_id, type_id, text))
        comment_id += 1

    if rows:
        execute_batch(cur,
            "INSERT INTO ProteinComment (accession, comment_id, commentType_id, comment_text) "
            "VALUES (%s, %s, %s, %s) ON CONFLICT DO NOTHING",
            rows)


def insert_features(cur, entry, accession, feature_type_cache):
    """Insert all features for a protein."""
    features = entry.get("features", [])
    feature_id = 1
    rows = []
    for feature in features:
        ftype = feature.get("type", "Unknown")
        description = feature.get("description", "")
        start, end = get_feature_location(feature)

        # Resolve or create FeatureType
        if ftype not in feature_type_cache:
            cur.execute(
                "INSERT INTO FeatureType (type_name) VALUES (%s) RETURNING type_id",
                (ftype,)
            )
            feature_type_cache[ftype] = cur.fetchone()[0]

        type_id = feature_type_cache[ftype]
        rows.append((accession, feature_id, type_id, description, start, end))
        feature_id += 1

    if rows:
        execute_batch(cur,
            "INSERT INTO ProteinFeature "
            "(accession, feature_id, featureType_id, description, start_position, end_position) "
            "VALUES (%s, %s, %s, %s, %s, %s) ON CONFLICT DO NOTHING",
            rows)


def insert_go_terms(cur, entry, accession):
    """Insert GO terms and protein-go associations."""
    go_annotations = entry.get("goAnnotations", [])
    go_rows = []
    assoc_rows = []
    for go in go_annotations:
        go_id = go.get("id")
        if not go_id:
            continue
        term_name = go.get("term", "")
        aspect = go.get("aspect", "P")
        category = ASPECT_MAP.get(aspect, "")

        go_rows.append((go_id, term_name, category, aspect))
        assoc_rows.append((accession, go_id))

    if go_rows:
        execute_batch(cur,
            "INSERT INTO GoTerm (go_id, term_name, category, aspect) "
            "VALUES (%s, %s, %s, %s) ON CONFLICT (go_id) DO NOTHING",
            go_rows)
    if assoc_rows:
        execute_batch(cur,
            "INSERT INTO annotated_with (accession, go_id) "
            "VALUES (%s, %s) ON CONFLICT DO NOTHING",
            assoc_rows)


def insert_keywords(cur, entry, accession):
    """Insert keywords and protein-keyword associations."""
    keywords = entry.get("keywords", [])
    kw_rows = []
    assoc_rows = []
    for kw in keywords:
        kw_id = kw.get("id")
        kw_name = kw.get("name", "")
        if not kw_id:
            continue
        kw_rows.append((kw_id, kw_name))
        assoc_rows.append((accession, kw_id))

    if kw_rows:
        execute_batch(cur,
            "INSERT INTO Keyword (keyword_id, keyword_name) "
            "VALUES (%s, %s) ON CONFLICT (keyword_id) DO NOTHING",
            kw_rows)
    if assoc_rows:
        execute_batch(cur,
            "INSERT INTO tagged_with (accession, keyword_id) "
            "VALUES (%s, %s) ON CONFLICT DO NOTHING",
            assoc_rows)


def insert_cross_references(cur, entry, accession):
    """Insert cross-references to external databases."""
    xrefs = entry.get("uniProtKBCrossReferences", [])
    rows = []
    for xref in xrefs:
        db_name = xref.get("database", "")
        ref_id = xref.get("id", "")
        if not db_name or not ref_id:
            continue
        rows.append((accession, ref_id, db_name))

    if rows:
        execute_batch(cur,
            "INSERT INTO CrossReference (accession, reference_id, database_name) "
            "VALUES (%s, %s, %s) ON CONFLICT DO NOTHING",
            rows)


# ============================================================
# Main import logic
# ============================================================

def import_proteins():
    """Main function: search proteins, fetch full entries, insert into DB."""
    print(f"Searching UniProt for proteins: {QUERY}")
    accessions = search_proteins(QUERY, PROTEIN_COUNT)
    print(f"Found {len(accessions)} proteins to import")

    if not accessions:
        print("No proteins found. Check your query.")
        sys.exit(1)

    # Connect to PostgreSQL
    conn = psycopg2.connect(DB_URL)
    conn.autocommit = False
    cur = conn.cursor()

    # Caches for lookup tables (avoid repeated SELECTs)
    comment_type_cache = {}
    feature_type_cache = {}

    success = 0
    errors = 0

    for i, accession in enumerate(accessions, 1):
        try:
            print(f"[{i}/{len(accessions)}] Fetching {accession}...", end=" ")
            entry = fetch_protein(accession)

            # Insert in dependency order:
            # 1. Organism (parent of Protein)
            insert_organism(cur, entry)

            # 2. Protein (parent of all detail tables)
            acc = insert_protein(cur, entry)
            if not acc:
                print("SKIP (already exists)")
                continue

            # 3. Lookup tables (GoTerm, Keyword) and associations
            insert_go_terms(cur, entry, accession)
            insert_keywords(cur, entry, accession)

            # 4. Detail tables (depend on Protein + lookup tables)
            insert_comments(cur, entry, accession, comment_type_cache)
            insert_features(cur, entry, accession, feature_type_cache)
            insert_cross_references(cur, entry, accession)

            conn.commit()
            success += 1
            print("OK")

            # Rate limiting: be nice to the API
            time.sleep(0.5)

        except requests.HTTPError as e:
            conn.rollback()
            print(f"HTTP ERROR: {e}")
            errors += 1
        except psycopg2.Error as e:
            conn.rollback()
            print(f"DB ERROR: {e}")
            errors += 1
        except Exception as e:
            conn.rollback()
            print(f"ERROR: {e}")
            errors += 1

    # Print summary
    print(f"\n{'=' * 50}")
    print(f"Import complete: {success} proteins imported, {errors} errors")

    # Print table counts
    print(f"\n{'=' * 50}")
    print("Table row counts:")
    tables = [
        "Organism", "Protein", "GoTerm", "Keyword",
        "CommentType", "FeatureType",
        "ProteinComment", "ProteinFeature", "CrossReference",
        "annotated_with", "tagged_with"
    ]
    for table in tables:
        cur.execute(f"SELECT COUNT(*) FROM {table}")
        count = cur.fetchone()[0]
        print(f"  {table}: {count} rows")

    cur.close()
    conn.close()


if __name__ == "__main__":
    import_proteins()
