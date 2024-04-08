package account

import (
	"net/http"
	"testing"
	"time"

	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/test"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.Account{
		{"123", "person@example.com", "xyz", time.Now(), time.Now()},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/accounts", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/accounts/123", "", nil, http.StatusOK, `*{"id":"123","email":"person@example.com","firebase_id":"xyz"}*`},
		{"get unknown", "GET", "/accounts/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/accounts", `{"name":"test"}`, header, http.StatusCreated, "*test*"},
		{"create ok count", "GET", "/accounts", "", nil, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/accounts", `{"name":"test"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/accounts", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/accounts/123", `{"name":"accountxyz"}`, header, http.StatusOK, "*accountxyz*"},
		{"update verify", "GET", "/accounts/123", "", nil, http.StatusOK, `*accountxyz*`},
		{"update auth error", "PUT", "/accounts/123", `{"name":"accountxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/accounts/123", `"name":"accountxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/accounts/123", ``, header, http.StatusOK, "*accountxyz*"},
		{"delete verify", "DELETE", "/accounts/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/accounts/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
