package common

const (
	AppEnvDev = "dev"
)

const InviteTokenExpirySecond = 604800

const CurrentUser = "user"

type Requester interface {
	GetRole() string
}
