package mock

import (
	"app-invite-service/common"
	"app-invite-service/component/tokenprovider"
	"app-invite-service/modules/user/usermodel"
	"context"
	"time"
)

type mockUserStore struct{}

func NewMockUserStore() *mockUserStore {
	return &mockUserStore{}
}

func (m *mockUserStore) FindUser(_ context.Context, conditions map[string]interface{}, _ ...string) (*usermodel.User, error) {
	if val, ok := conditions["email"]; ok && val.(string) == "user@gmail.com" {
		return &usermodel.User{Email: val.(string), Password: "user@123", Status: 1, Salt: ""}, nil
	}
	if val, ok := conditions["id"]; ok && val.(int) == 1 {
		return &usermodel.User{Email: "user@gmail.com", Password: "user@123", Status: 1, Salt: ""}, nil
	}
	if val, ok := conditions["id"]; ok && val.(int) == 2 {
		return &usermodel.User{Email: "user2@gmail.com", Password: "user2@123", Status: 1, Salt: ""}, nil
	}
	return nil, common.ErrRecordNotFound
}

func (m *mockUserStore) CreateUser(_ context.Context, data *usermodel.UserCreate) error {
	data.Id = 3
	return nil
}

type mockProvider struct{}

func NewMockProvider() *mockProvider {
	return &mockProvider{}
}

func (m *mockProvider) Generate(_ tokenprovider.TokenPayload, expiry int) (*tokenprovider.Token, error) {
	return &tokenprovider.Token{
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7InVzZXJfaWQiOjF9LCJleHAiOjE2NTM0MDk1MDksImlhdCI6MTY1MzMyMzEwOX0.SYzR9JXyIc_VeuLXLAnxFWTM3nO6LQfWbyO-vTK3fMo",
		Expiry:  expiry,
		Created: time.Now().UTC(),
	}, nil
}

func (m *mockProvider) Validate(_ string) (*tokenprovider.TokenPayload, error) {
	return &tokenprovider.TokenPayload{}, nil
}

type mockHash struct{}

func NewMockHash() *mockHash {
	return &mockHash{}
}

func (m *mockHash) Hash(data string) string {
	return data
}
