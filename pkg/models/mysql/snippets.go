package mysql

import (
	"database/sql"
	"errors"

	"github.com/eiliz/snippetbox/pkg/models"
)

// SnippetModel is a type that wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a new snippet into the db
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// Use backticks to spread statement in multiple lines
	// DB.Exec does 3 steps: creates a prepared statement which the database
	// parses, compiles and stores for execution; passes the parameter values to
	// the database and the statement gets executed - because the parameters are
	// passed after the statement has been compiled they are treated as just data,
	// they cannot result into an SQL injection; finally the prepared statement is
	// closed/deallocated.
	stmt := `INSERT INTO snippets (title, content, created, expires)
					VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// The LastInsertId method of the result object returns the ID our newly
	// inserted record in the snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	// The ID returned is of type int64, we need to convert to int before returning
	return int(id), nil
}

// Get returns a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
					WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

// Latest returns the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
					WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	snippets := []*models.Snippet{}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
