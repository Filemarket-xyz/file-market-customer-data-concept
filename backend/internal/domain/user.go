package domain

import "github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"

type (
	Role int
)

const (
	RoleClient = Role(iota)
)

func RoleToModel(role Role) models.UserRole {
	switch role {
	case RoleClient:
		return models.UserRoleClient
	}
	return ""
}

type UserWithTokenNumber struct {
	Id         int64
	Role       Role
	Number     int64
	Authorized bool
}

type AuthMessage struct {
	Address   string
	Message   string
	CreatedAt int64
}
