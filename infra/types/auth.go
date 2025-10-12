package types

import "subscriptionplus/server/infra/store/postgres/models"

type AuthToken struct {
	User models.UserInfo `json:"user"`
}
