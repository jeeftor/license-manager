package language

import (
	"github.com/jeeftor/license-manager/internal/logger"
	"github.com/jeeftor/license-manager/internal/styles"
)

type JavaHandler struct {
	*GenericHandler
}

func NewJavaHandler(logger *logger.Logger, style styles.HeaderFooterStyle) *JavaHandler {
	h := &JavaHandler{GenericHandler: NewGenericHandler(logger, style, "cpp")}
	h.GenericHandler.subclassHandler = h
	return h
}
