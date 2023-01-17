package hash_test

import (
	"app-invite-service/component/hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMd5Hash_Hash(t *testing.T) {
	var tcs = []struct {
		password string
		salt     string
		expected string
	}{
		{"nana@123", "BMOcrdlEltpGCZxZkmVyBqyDwxrDXkxPLZMOFDSXNxGqrwKoxt", "b0dd9c5cfd02c3e96171ab3f08e67dac"},
		{"root@123", "BMOcrdlEltpGCZxZkmVyBqyDwxrDXkxPLZMOFDSXNxGqrwKoxt", "148a188a95d7ce98c5e45badb5a39f3f"},
		{"password!456", "tNCdTDEnAVLkqXKcyOEpAgEsPkTKhEgzxKRGyvZomTqlkrzwxR", "e6ce76f34dda7b86a51a1aa8cd4ac040"},
	}

	for _, tc := range tcs {
		output := hash.NewMd5Hash().Hash(tc.password + tc.salt)
		assert.Equal(t, tc.expected, output, "they should be equal")
	}
}
