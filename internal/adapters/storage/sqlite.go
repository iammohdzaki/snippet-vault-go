package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"snippet-vault-go/internal/core"
	"strconv"

	_ "modernc.org/sqlite"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(dbPath string) (*SQLiteRepo, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS snippets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		language TEXT NOT NULL,
		code TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	CREATE TABLE IF NOT EXISTS links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return &SQLiteRepo{db: db}, nil
}

func (r *SQLiteRepo) Save(snippet *core.Snippet) error {
	res, err := r.db.Exec(
		`INSERT INTO snippets (title, language, code, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		snippet.Title,
		snippet.Language,
		snippet.Code,
		snippet.CreatedAt,
		snippet.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	snippet.ID = strconv.FormatInt(id, 10)
	return nil
}

func (r *SQLiteRepo) GetAll() ([]core.Snippet, error) {
	rows, err := r.db.Query(`
		SELECT id, title, language, code, created_at, updated_at 
		FROM snippets 
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []core.Snippet
	for rows.Next() {
		var s core.Snippet
		if err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Language,
			&s.Code,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if snippets == nil {
		return []core.Snippet{}, nil
	}
	return snippets, nil
}

func (r *SQLiteRepo) Update(snippet *core.Snippet) error {
	res, err := r.db.Exec(
		`UPDATE snippets 
		 SET title = ?, language = ?, code = ?, updated_at = ? 
		 WHERE id = ?`,
		snippet.Title,
		snippet.Language,
		snippet.Code,
		snippet.UpdatedAt,
		snippet.ID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return core.ErrNotFound
	}

	// Populate the original created_at timestamp onto the struct
	_ = r.db.QueryRow("SELECT created_at FROM snippets WHERE id = ?", snippet.ID).Scan(&snippet.CreatedAt)
	return nil
}

func (r *SQLiteRepo) Delete(id string) error {
	res, err := r.db.Exec("DELETE FROM snippets WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return core.ErrNotFound
	}
	return nil
}

func (r *SQLiteRepo) SaveLink(link *core.Link) error {
	res, err := r.db.Exec(
		`INSERT INTO links (title, url, description, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		link.Title,
		link.URL,
		link.Description,
		link.CreatedAt,
		link.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	link.ID = strconv.FormatInt(id, 10)
	return nil
}

func (r *SQLiteRepo) GetAllLinks() ([]core.Link, error) {
	rows, err := r.db.Query(`
		SELECT id, title, url, description, created_at, updated_at 
		FROM links 
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []core.Link
	for rows.Next() {
		var l core.Link
		if err := rows.Scan(&l.ID, &l.Title, &l.URL, &l.Description, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	if links == nil {
		return []core.Link{}, nil
	}
	return links, nil
}

func (r *SQLiteRepo) UpdateLink(link *core.Link) error {
	res, err := r.db.Exec(
		`UPDATE links 
		 SET title = ?, url = ?, description = ?, updated_at = ? 
		 WHERE id = ?`,
		link.Title,
		link.URL,
		link.Description,
		link.UpdatedAt,
		link.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return core.ErrNotFound
	}
	_ = r.db.QueryRow("SELECT created_at FROM links WHERE id = ?", link.ID).Scan(&link.CreatedAt)
	return nil
}

func (r *SQLiteRepo) DeleteLink(id string) error {
	res, err := r.db.Exec("DELETE FROM links WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return core.ErrNotFound
	}
	return nil
}
