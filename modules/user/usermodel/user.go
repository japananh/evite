package usermodel

import (
	"app-invite-service/common"
	"app-invite-service/component/tokenprovider"
	"errors"
	"strings"
	"time"
	"unicode"
)

const EntityName = "User"

var (
	ErrInvalidCharacterMsg         = "password has invalid characters"
	ErrNotEnoughCharacterMsg       = "password must have at least 8 characters"
	ErrMustHaveNumberMsg           = "password must have at least 1 number"
	ErrMustHaveLetterMsg           = "password must have at least 1 letter"
	ErrMustHaveSpecialCharacterMsg = "password must have at least 1 special character"
)

var ErrEmailOrPasswordInvalid = common.NewCustomError(
	errors.New("email or password invalid"),
	"email or password invalid",
	"ErrEmailOrPasswordInvalid",
)

func ErrPasswordInvalid(msg string) *common.AppError {
	return common.NewCustomError(
		errors.New(msg),
		msg,
		"ErrPasswordInvalid",
	)
}

type User struct {
	Id        int        `json:"-" gorm:"column:id;"`
	Status    int        `json:"status" gorm:"column:status;default:1;"`
	Email     string     `json:"email" form:"email" binding:"required" gorm:"column:email;"`
	Password  string     `json:"password" form:"password" binding:"required" gorm:"column:password;"`
	Salt      string     `json:"-" gorm:"column:salt;"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) GetUserId() int {
	return u.Id
}

func (u *User) GetEmail() string {
	return u.Email
}

type UserCreate struct {
	Id        int        `json:"-" gorm:"column:id;"`
	Status    int        `json:"status" gorm:"column:status;default:1;"`
	Email     string     `json:"email" form:"email" binding:"required" gorm:"column:email;"`
	Password  string     `json:"password" form:"password" binding:"required" gorm:"column:password;"`
	Salt      string     `json:"-" gorm:"column:salt;"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

func (u *UserCreate) Validate() error {
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)

	if errMsg := VerifyPassword(u.Password); errMsg != "" {
		return ErrPasswordInvalid(errMsg)
	}

	return nil
}

func VerifyPassword(s string) string {
	hasNumber, hasLetter, hasSpecial, hasInvalidCharacter := false, false, false, false
	letterCount := 0

	for _, c := range s {
		letterCount++
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		default:
			hasInvalidCharacter = true
		}
	}

	if hasInvalidCharacter {
		return ErrInvalidCharacterMsg
	}

	if letterCount < 8 {
		return ErrNotEnoughCharacterMsg
	}

	if !hasNumber {
		return ErrMustHaveNumberMsg
	}

	if !hasLetter {
		return ErrMustHaveLetterMsg
	}

	if !hasSpecial {
		return ErrMustHaveSpecialCharacterMsg
	}

	return ""
}

type UserLogin struct {
	Email    string `json:"email" form:"email" binding:"required" gorm:"column:email;"`
	Password string `json:"password" form:"password" binding:"required" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

type UserLoginWithInviteToken struct {
	InvitationToken string `json:"invite_token" form:"invite_token" binding:"required"`
}

func (u *UserLoginWithInviteToken) Validate() error {
	u.InvitationToken = strings.TrimSpace(u.InvitationToken)
	return nil
}

type Account struct {
	AccessToken  *tokenprovider.Token `json:"access_token"`
	RefreshToken *tokenprovider.Token `json:"refresh_token"`
}

func NewAccount(at, rt *tokenprovider.Token) *Account {
	return &Account{
		AccessToken:  at,
		RefreshToken: rt,
	}
}

type InvitationToken struct {
	Token string `json:"invite_token"`
}
