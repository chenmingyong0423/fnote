// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package visit_log

import (
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/internal/web"
)

type (
	Handler      = web.VisitLogHandler
	Service      = service.IVisitLogService
	TendencyData = domain.TendencyData
	Module       struct {
		Svc Service
		Hdl *Handler
	}
)
