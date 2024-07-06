package data

import (
	"database/sql"
	"errors"
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

// The Insert() method accepts a pointer to a movie struct, which should contain the
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

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the snip struct.
	return m.DB.QueryRow(query, args...).Scan(&snip.ID, &snip.CreatedAt, &snip.Version)
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

	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Snip struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRow(query, id).Scan(
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

// Add a placeholder method for updating a specific record in the movies table.
func (m SnipModel) Update(movie *Snip) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m SnipModel) Delete(id int64) error {
	return nil
}
