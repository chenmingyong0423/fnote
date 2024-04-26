package data_analysis

import "github.com/chenmingyong0423/fnote/server/internal/data_analysis/internal/web"

type (
	Handler = web.DataAnalysisHandler
)
type Module struct {
	Hdl *Handler
}
