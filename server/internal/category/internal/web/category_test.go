package web

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/category/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/category/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type routeCategoryService struct {
	service.ICategoryService
	err error
}

func (s routeCategoryService) GetCategoryByRoute(context.Context, string) (domain.Category, error) {
	return domain.Category{Name: "Backend"}, s.err
}

func TestCategoryRouteStatus(t *testing.T) {
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
			NewCategoryHandler(routeCategoryService{err: tc.err}).RegisterGinRoutes(router)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest("GET", "/categories/route/backend", nil))
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
		})
	}
}
