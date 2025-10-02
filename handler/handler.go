package handler

import (
	"database/sql"
	"subscriptionplus/server/infra/logger"
	"subscriptionplus/server/infra/store/postgres/store"
)

type BaseHandler struct {
	Db     *sql.DB
	Logger *logger.Logger
	Store  store.Storage
}
