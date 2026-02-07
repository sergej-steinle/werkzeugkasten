package werkzeugkasten

import (
	"math/rand/v2"
	"strings"
)

const randomStrSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+_"

// Werkzeug provides utility functions.
type Werkzeug struct {
}

// RandomString returns a random string of length n using characters from randomStrSource.
func (t *Werkzeug) RandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(randomStrSource[rand.IntN(len(randomStrSource))])
	}
	return sb.String()
}
