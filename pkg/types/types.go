package types

// Rate represents the exchange rate between two currencies.
// It is expressed as a float32 value.
type Rate float32

// User represents an email address of a subscriber.
// It is expressed as a string value.
type User struct {
	ID    int
	Email string
}
