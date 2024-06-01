package handlers

import (
	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
	Handler
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

// HelloWorldHandle return a simple hello world message
//
//	@Summary		Return a simple hello world message
//	@Description	Return a simple hello world message
//	@Tags			v1
//	@Produce		plain
//	@Success		200	{object}	string
//	@Router			/v1/hello [get]
func (h *IndexHandler) HelloWorldHandle(ctx echo.Context) error {
	return ctx.String(200, "Hello, World!")
}
