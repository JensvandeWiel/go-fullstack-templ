package handlers

import (
	"github.com/JensvandeWiel/middleware"
	"github.com/labstack/echo/v4"
)

type IndexInertiaHandler struct {
	Handler
}

func NewIndexInertiaHandler() *IndexInertiaHandler {
	return &IndexInertiaHandler{}
}

// IndexInertiaHandle return a simple hello world message
func (h *IndexInertiaHandler) IndexInertiaHandle(ctx echo.Context) error {
	i := middleware.GetInertia(ctx)
	return i.Render(ctx.Response(), ctx.Request(), "Index", nil)
}
