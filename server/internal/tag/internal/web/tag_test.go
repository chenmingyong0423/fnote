package web

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type routeTagService struct {
	service.ITagService
	err error
}

func (s routeTagService) GetTagByRoute(context.Context, string) (domain.Tag, error) {
	return domain.Tag{Name: "Go"}, s.err
}

func TestTagRouteStatus(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"existing", nil, 200},
		{"missing", fmt.Errorf("lookup: %w", mongo.ErrNoDocuments), 404},
		{"unavailable", errors.New("database unavailable"), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			NewTagHandler(routeTagService{err: tc.err}).RegisterGinRoutes(router)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest("GET", "/tags/route/go", nil))
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
		})
	}
}
