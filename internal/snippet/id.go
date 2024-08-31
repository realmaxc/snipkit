package snippet

import (
	"crypto/rand"
	"encoding/base32"
)

// idEncoding is a Crockford-flavoured base32 alphabet: lowercase, no padding,
// and no visually ambiguous characters (i, l, o, u).
var idEncoding = base32.NewEncoding("0123456789abcdefghjkmnpqrstvwxyz").WithPadding(base32.NoPadding)

// NewID returns a short, unambiguous, lowercase identifier for a snippet.
// Five random bytes encode to eight base32 characters.
func NewID() string {
	var b [5]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand.Read only fails if the system entropy source is
		// unavailable, which we cannot sensibly recover from here.
		panic("snippet: reading random bytes: " + err.Error())
	}
	return idEncoding.EncodeToString(b[:])
}
