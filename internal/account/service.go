package account

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// Service encapsulates usecase logic for Accounts.
type Service interface {
	Get(ctx context.Context, id int, email string, firebaseId string) (Account, error)
	Query(ctx context.Context, offset, limit int) ([]Account, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateAccountRequest) (Account, error)
	Update(ctx context.Context, id int, email string, firebaseId string, input UpdateAccountRequest) (Account, error)
	Delete(ctx context.Context, id int, email string, firebaseId string) (Account, error)
}

// Account represents the data about an Account.
type Account struct {
	entity.Account
}

// CreateAccountRequest represents an Account creation request.
type CreateAccountRequest struct {
	Email      string `json:"email"`
	FirebaseId string `json:"FirebaseId"`
}

// Validate validates the CreateAccountRequest fields.
func (m CreateAccountRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, validation.Length(0, 128)),
	)
}

// UpdateAccountRequest represents an Account update request.
type UpdateAccountRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateAccountRequest fields.
func (m UpdateAccountRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new Account service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the Account with the specified the Account ID.
func (s service) Get(ctx context.Context, id int, email string, firebaseId string) (Account, error) {
	account, err := s.repo.Get(ctx, id, email, firebaseId)
	if err != nil {
		return Account{}, err
	}
	return Account{account}, nil
}

// Create creates a new Account.
func (s service) Create(ctx context.Context, req CreateAccountRequest) (Account, error) {
	if err := req.Validate(); err != nil {
		return Account{}, err
	}
	now := time.Now()
	err := s.repo.Create(ctx, entity.Account{
		Email:      req.Email,
		FirebaseId: req.FirebaseId,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return Account{}, err
	}
	return s.Get(ctx, 0, req.Email, "")
}

// Update updates the Account with the specified ID.
func (s service) Update(ctx context.Context, id int, email string, firebaseId string, req UpdateAccountRequest) (Account, error) {
	if err := req.Validate(); err != nil {
		return Account{}, err
	}

	Account, err := s.Get(ctx, id, email, firebaseId)
	if err != nil {
		return Account, err
	}
	Account.Email = req.Name
	Account.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, Account.Account); err != nil {
		return Account, err
	}
	return Account, nil
}

// Delete deletes the Account with the specified ID.
func (s service) Delete(ctx context.Context, id int, email string, firebaseId string) (Account, error) {
	account, err := s.Get(ctx, id, email, firebaseId)
	if err != nil {
		return Account{}, err
	}
	if err = s.repo.Delete(ctx, id, email, firebaseId); err != nil {
		return Account{}, err
	}
	return account, nil
}

// Count returns the number of Accounts.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the Accounts with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int) ([]Account, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []Account{}
	for _, item := range items {
		result = append(result, Account{item})
	}
	return result, nil
}
