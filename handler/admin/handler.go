package admin

import "subscriptionplus/server/handler"

type Handler struct {
	*handler.BaseHandler
}

func NewHandler(base *handler.BaseHandler) *Handler {
	return &Handler{BaseHandler: base}
}
