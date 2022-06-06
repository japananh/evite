package userbiz_test

import (
	"app-invite-service/component/tokenprovider"
	"app-invite-service/mock"
	"app-invite-service/modules/user/userbiz"
	"app-invite-service/modules/user/usermodel"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoginBiz_Login(t *testing.T) {
	tcs := []struct {
		atExpiry    int
		rtExpiry    int
		email       string
		password    string
		expectedErr error
	}{
		{8600, 60800, "user@gmail.com", "user@123", nil},
		{8600, 60800, "user@gmail.com", "user@1234", errors.New("email or password invalid")},
		{8600, 60800, "user@gmail.com", "", errors.New("email or password invalid")},
		{8600, 60800, "user1@gmail.com", "user@123", errors.New("email or password invalid")},
	}

	for _, tc := range tcs {
		biz := userbiz.NewLoginBiz(
			mock.NewMockUserStore(),
			mock.NewMockProvider(),
			mock.NewMockHash(),
			&tokenprovider.TokenConfig{AccessTokenExpiry: tc.atExpiry, RefreshTokenExpiry: tc.rtExpiry},
		)
		user, err := biz.Login(nil, &usermodel.UserLogin{Email: tc.email, Password: tc.password})
		if tc.expectedErr != nil {
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Equal(t, tc.expectedErr.Error(), err.Error())
		} else {
			assert.NotNil(t, user)
		}
	}
}
