package domain

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// Service encapsulates usecase logic for Domains.
type Service interface {
	Get(ctx context.Context, id int, accountId int) (Domain, error)
	Query(ctx context.Context, offset, limit int, accountId int) ([]Domain, error)
	Count(ctx context.Context, accountId int) (int, error)
	Create(ctx context.Context, input CreateDomainRequest) (Domain, error)
	Update(ctx context.Context, id int, input UpdateDomainRequest) (Domain, error)
	Delete(ctx context.Context, id int, accountId int) (Domain, error)
}

// Domain represents the data about an Domain.
type Domain struct {
	entity.Domain
}

// CreateDomainRequest represents an Domain creation request.
type CreateDomainRequest struct {
	Name      string `json:"name"`
	AccountId int    `json:"account_id"`
}

// Validate validates the CreateDomainRequest fields.
func (m CreateDomainRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.AccountId, validation.Required, validation.Min(0)),
	)
}

// UpdateDomainRequest represents an Domain update request.
type UpdateDomainRequest struct {
	Name      string `json:"name"`
	AccountId int    `json:"account_id"`
}

// Validate validates the CreateDomainRequest fields.
func (m UpdateDomainRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new Domain service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the Domain with the specified the Domain ID.
func (s service) Get(ctx context.Context, id int, accountId int) (Domain, error) {
	domain, err := s.repo.Get(ctx, id, accountId)
	if err != nil {
		return Domain{}, err
	}
	return Domain{domain}, nil
}

// Create creates a new Domain.
func (s service) Create(ctx context.Context, req CreateDomainRequest) (Domain, error) {
	if err := req.Validate(); err != nil {
		return Domain{}, err
	}
	now := time.Now()
	domain, err := s.repo.Create(ctx, entity.Domain{
		Domain:    req.Name,
		AccountId: req.AccountId,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return Domain{}, err
	}
	return s.Get(ctx, domain.ID, domain.AccountId)
}

// Update updates the Domain with the specified ID.
func (s service) Update(ctx context.Context, id int, req UpdateDomainRequest) (Domain, error) {
	if err := req.Validate(); err != nil {
		return Domain{}, err
	}

	Domain, err := s.Get(ctx, id, req.AccountId)
	if err != nil {
		return Domain, err
	}
	Domain.Domain.Domain = req.Name
	Domain.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, Domain.Domain); err != nil {
		return Domain, err
	}
	return Domain, nil
}

// Delete deletes the Domain with the specified ID.
func (s service) Delete(ctx context.Context, id int, accountId int) (Domain, error) {
	domain, err := s.Get(ctx, id, accountId)
	if err != nil {
		return Domain{}, err
	}
	if err = s.repo.Delete(ctx, id, accountId); err != nil {
		return Domain{}, err
	}
	return domain, nil
}

// Count returns the number of Domains.
func (s service) Count(ctx context.Context, accountId int) (int, error) {
	return s.repo.Count(ctx, accountId)
}

// Query returns the Domains with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int, accountId int) ([]Domain, error) {
	items, err := s.repo.Query(ctx, offset, limit, accountId)
	if err != nil {
		return nil, err
	}
	result := []Domain{}
	for _, item := range items {
		result = append(result, Domain{item})
	}
	return result, nil
}
