package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pwilliams-ck/sniplate/internal/validator"
)

type Snip struct {
	ID        int64     `json:"id"`                // Unique integer ID for the snip
	CreatedAt time.Time `json:"created_at"`        // Timestamp for when the snip is added to our database
	Title     string    `json:"title"`             // Snip title
	Content   string    `json:"content,omitempty"` // Content of the snip
	Tags      []string  `json:"tags,omitempty"`    // Slice of tags for the snip
	Version   int32     `json:"version"`           // Starts at 1 and increments each time the snip is updated
}

func ValidateSnip(v *validator.Validator, snip *Snip) {
	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true.  In the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(snip.Title != "", "title", "must be provided")
	v.Check(len(snip.Title) <= 1000, "title", "must not be more than 500 bytes long")

	v.Check(len(snip.Content) <= 100000, "content", "must not be more than 100000 bytes long")

	v.Check(len(snip.Tags) >= 0, "tags", "may not contain negative tags")
	v.Check(len(snip.Tags) <= 10, "tags", "must contain no more than 5 tags")
	v.Check(validator.Unique(snip.Tags), "tags", "must not contain duplicate values")

	// Check that each tag is not an empty string.
	for _, tag := range snip.Tags {
		v.Check(tag != "", "tags", "tags must not contain empty values")
	}
}

// Define a SnipModel struct type which wraps a sql.DB connection pool.
type SnipModel struct {
	DB *sql.DB
}

// The Insert() method accepts a pointer to a snip struct, which should contain the
// data for the new record.
func (m SnipModel) Insert(snip *Snip) error {
	// Define the SQL query for inserting a new record in the snips table and returning
	// the system-generated data.
	query := `
        INSERT INTO snips (title, content, tags)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, version`

	// Create an args slice containing the values for the placeholder parameters from
	// the snip struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{snip.Title, snip.Content, pq.Array(snip.Tags)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the snip struct.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&snip.ID, &snip.CreatedAt, &snip.Version)
}

func (m SnipModel) Get(id int64) (*Snip, error) {
	// The PostgreSQL bigserial type that we're using for the snip ID starts
	// auto-incrementing at 1 by default, so we know that no snips will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Define the SQL query for retrieving the snip data.
	query := `
        SELECT id, created_at, title, content, tags, version
        FROM snips
        WHERE id = $1`

	// Declare a Snip struct to hold the data returned by the query.
	var snip Snip

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Snip struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&snip.ID,
		&snip.CreatedAt,
		&snip.Title,
		&snip.Content,
		pq.Array(&snip.Tags),
		&snip.Version,
	)
	// Handle any errors. If there was no matching snip found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Otherwise, return a pointer to the Snip struct.
	return &snip, nil
}

// GetAll() returns a slice of snips. Although we're not
// using them right now, we've set this up to accept the
// various filter parameters as arguments.
func (m SnipModel) GetAll(title string, tags []string, filters Filters) ([]*Snip, error) {
	// Construct the SQL query to retrieve all snip records.
	query := fmt.Sprintf(`
        SELECT id, created_at, title, content, tags, version
        FROM snips
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (tags @> $2 OR $2 = '{}')
        ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(tags))
	if err != nil {
		return nil, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()

	// Initialize an empty slice to hold the snip data.
	snips := []*Snip{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Snip struct to hold the data for an individual snip.
		var snip Snip

		// Scan the values from the row into the Snip struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&snip.ID,
			&snip.CreatedAt,
			&snip.Title,
			&snip.Content,
			pq.Array(&snip.Tags),
			&snip.Version,
		)
		if err != nil {
			return nil, err
		}

		// Add the Snip struct to the slice.
		snips = append(snips, &snip)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK, then return the slice of snips.
	return snips, nil
}

func (m SnipModel) Update(snip *Snip) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
        UPDATE snips 
        SET title = $1, content = $2, tags = $3, version = version + 1
        WHERE id = $4 AND version = $5
        RETURNING version`

	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		snip.Title,
		snip.Content,
		pq.Array(snip.Tags),
		snip.ID,
		snip.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&snip.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m SnipModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the snip ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	// Construct the SQL query to delete the record.
	query := `
        DELETE FROM snips
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, we know that the snips table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
