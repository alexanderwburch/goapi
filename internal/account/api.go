package account

import (
	"net/http"

	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/qiangxue/go-rest-api/pkg/pagination"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/accounts/<id>", res.get)
	r.Get("/accounts", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/accounts", res.create)
	r.Put("/accounts/<id>", res.update)
	r.Delete("/accounts/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	account, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(account)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	accounts, err := r.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = accounts
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateAccountRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	account, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(account, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateAccountRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	account, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(account)
}

func (r resource) delete(c *routing.Context) error {
	account, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(account)
}
