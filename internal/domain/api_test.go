package domain

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
	repo := &mockRepository{items: []entity.Domain{
		{"123", "12345", "example.com", time.Now(), time.Now()},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/domains", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/domains/123", "", nil, http.StatusOK, `*{"id":"123","account_id":"12345","domain":"example.com"}*`},
		{"get unknown", "GET", "/domains/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/domains", `{"name":"test.com","account_id":"12345"}`, header, http.StatusCreated, "*test.com*"},
		{"create ok count", "GET", "/domains", "", nil, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/domains", `{"name":"test"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/domains", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/domains/123", `{"name":"domainxyz"}`, header, http.StatusOK, "*domainxyz*"},
		{"update verify", "GET", "/domains/123", "", nil, http.StatusOK, `*domainxyz*`},
		{"update auth error", "PUT", "/domains/123", `{"name":"domainxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/domains/123", `"name":"domainxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/domains/123", ``, header, http.StatusOK, "*domainxyz*"},
		{"delete verify", "DELETE", "/domains/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/domains/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
