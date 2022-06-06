package common

const (
	DbTypeUser = 1
	DbTypeTask = 2
)

const InviteTokenExpirySecond = 604800

const CurrentUser = "user"

type Requester interface {
	GetUserId() int
	GetEmail() string
}
