package domain

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/qiangxue/go-rest-api/pkg/pagination"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/domains", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/domains", res.create)
	r.Put("/domains/<id>", res.update)
	r.Delete("/domains", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	accountId, err := strconv.Atoi(c.Param("account_id"))
	domain, err := r.service.Get(c.Request.Context(), c.Param("domain"), accountId)
	if err != nil {
		return err
	}

	return c.Write(domain)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	accountId, err := strconv.Atoi(c.Query("account_id"))
	count, err := r.service.Count(ctx, accountId)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	domains, err := r.service.Query(ctx, pages.Offset(), pages.Limit(), accountId)
	if err != nil {
		return err
	}
	pages.Items = domains
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateDomainRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	domain, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(domain, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateDomainRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	domain, err := r.service.Update(c.Request.Context(), c.Param("domain"), input)
	if err != nil {
		return err
	}

	return c.Write(domain)
}

func (r resource) delete(c *routing.Context) error {
	accountId, err := strconv.Atoi(c.Query("account_id"))
	domainParam := c.Query("domain")
	print(fmt.Sprintf("api.go Deleting domain %s %s", domainParam, c.Query("account_id")))
	domain, err := r.service.Delete(c.Request.Context(), domainParam, accountId)
	if err != nil {
		return err
	}

	return c.Write(domain)
}
