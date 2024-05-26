package data

import (
	"time"
)

type Snip struct {
	ID        int64     `json:"id"`                // Unique integer ID for the snip
	CreatedAt time.Time `json:"created_at"`        // Timestamp for when the snip is added to our database
	Title     string    `json:"title"`             // Snip title
	Content   string    `json:"content,omitempty"` // Content of the snip
	Tags      []string  `json:"tags,omitempty"`    // Slice of tags for the snip
	Version   int32     `json:"version"`           // Starts at 1 and increments each time the snip is updated
}
