package common

const (
	DbTypeUser = 1
)

const InviteTokenExpirySecond = 604800

const CurrentUser = "user"

type Requester interface {
	GetRole() string
}
