package data

import (
	"time"

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
