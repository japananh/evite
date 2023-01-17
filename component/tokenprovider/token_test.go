package tokenprovider_test

import (
	"app-invite-service/component/tokenprovider"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToken_NewTokenConfig(t *testing.T) {
	var tcs = []struct {
		ate      int
		rte      int
		expected *tokenprovider.TokenConfig
	}{
		{84600, 604800, &tokenprovider.TokenConfig{AccessTokenExpiry: 84600, RefreshTokenExpiry: 604800}},
		{300, 6048, &tokenprovider.TokenConfig{AccessTokenExpiry: 300, RefreshTokenExpiry: 6048}},
		{200, 3000, &tokenprovider.TokenConfig{AccessTokenExpiry: 200, RefreshTokenExpiry: 3000}},
	}

	for _, tc := range tcs {
		output, err := tokenprovider.NewTokenConfig(tc.ate, tc.rte)
		require.Nil(t, err)
		assert.Equal(t, tc.expected, output, "they should be equal")
	}
}
