package domain

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateDomainRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateDomainRequest
		wantError bool
	}{
		{"success", CreateDomainRequest{Name: "test.com", AccountId: "1234"}, false},
		{"required", CreateDomainRequest{Name: ""}, true},
		{"too long", CreateDomainRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateDomainRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdateDomainRequest
		wantError bool
	}{
		{"success", UpdateDomainRequest{Name: "test"}, false},
		{"required", UpdateDomainRequest{Name: ""}, true},
		{"too long", UpdateDomainRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// successful creation
	domain, err := s.Create(ctx, CreateDomainRequest{Name: "example.com", AccountId: "1234"})
	assert.Nil(t, err)
	assert.NotEmpty(t, domain.ID)
	id := domain.ID
	assert.Equal(t, "example.com", domain.Domain.Domain)
	assert.NotEmpty(t, domain.CreatedAt)
	assert.NotEmpty(t, domain.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateDomainRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateDomainRequest{Name: "error", AccountId: "1234"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateDomainRequest{Name: "example.com"})

	// update
	domain, err = s.Update(ctx, id, UpdateDomainRequest{Name: "example.com updated"})
	assert.Nil(t, err)
	assert.Equal(t, "example.com updated", domain.Domain.Domain)
	_, err = s.Update(ctx, "none", UpdateDomainRequest{Name: "example.com"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdateDomainRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateDomainRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	domain, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "example.com updated", domain.Domain.Domain)
	assert.Equal(t, id, domain.ID)

	// query
	domains, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 1, len(domains))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	domain, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, domain.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 0, count)
}

type mockRepository struct {
	items []entity.Domain
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Domain, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.Domain{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.Domain, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, domain entity.Domain) error {
	if domain.Domain == "error" {
		return errCRUD
	}
	m.items = append(m.items, domain)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, domain entity.Domain) error {
	if domain.Domain == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == domain.ID {
			m.items[i] = domain
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
