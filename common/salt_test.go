package common_test

import (
	"app-invite-service/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSalt_GenSalt(t *testing.T) {
	var tcs = []struct {
		arg      int
		expected string
	}{
		{0, ""},
		{-1, "BMOcrdlEltpGCZxZkmVyBqyDwxrDXkxPLZMOFDSXNxGqrwKoxt"},
		{50, "BMOcrdlEltpGCZxZkmVyBqyDwxrDXkxPLZMOFDSXNxGqrwKoxt"},
		{20, "tNCdTDEnAVLkqXKcyOEp"},
	}

	for _, tc := range tcs {
		output := common.GenSalt(tc.arg)
		assert.Equal(t, len(output), len(tc.expected), "they should be equal")
	}
}
