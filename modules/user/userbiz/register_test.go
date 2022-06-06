package userbiz_test

import (
	"app-invite-service/mock"
	"app-invite-service/modules/user/userbiz"
	"app-invite-service/modules/user/usermodel"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserBiz_Register(t *testing.T) {
	tcs := []struct {
		email       string
		password    string
		expectedErr error
	}{
		{"user1@gmail.com", "user@123", nil},
		{"user@gmail.com", "user@123", errors.New("user already exists")},
		{"user2@gmail.com", "", errors.New("password must have at least 8 characters")},
		{"user3@gmail.com", "user", errors.New("password must have at least 8 characters")},
		{"user4@gmail.com", "password", errors.New("password must have at least 1 number")},
		{"user5@gmail.com", "12345678", errors.New("password must have at least 1 letter")},
		{"user6@gmail.com", "pass1234", errors.New("password must have at least 1 special character")},
		{"user7@gmail.com", "!@#$%^&*", errors.New("password must have at least 1 number")},
	}

	for _, tc := range tcs {
		biz := userbiz.NewRegisterBiz(
			mock.NewMockUserStore(),
			mock.NewMockHash(),
		)
		err := biz.Register(nil, &usermodel.UserCreate{Email: tc.email, Password: tc.password})
		if tc.expectedErr != nil {
			assert.Error(t, err)
			assert.Equal(t, tc.expectedErr.Error(), err.Error())
		}
	}
}
