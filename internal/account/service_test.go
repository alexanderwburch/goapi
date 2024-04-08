package account

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

func TestCreateAccountRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateAccountRequest
		wantError bool
	}{
		{"success", CreateAccountRequest{Name: "test"}, false},
		{"required", CreateAccountRequest{Name: ""}, true},
		{"too long", CreateAccountRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateAccountRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdateAccountRequest
		wantError bool
	}{
		{"success", UpdateAccountRequest{Name: "test"}, false},
		{"required", UpdateAccountRequest{Name: ""}, true},
		{"too long", UpdateAccountRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
	account, err := s.Create(ctx, CreateAccountRequest{Name: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, account.ID)
	id := account.ID
	assert.Equal(t, "test", account.Email)
	assert.NotEmpty(t, account.CreatedAt)
	assert.NotEmpty(t, account.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateAccountRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateAccountRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateAccountRequest{Name: "test2"})

	// update
	account, err = s.Update(ctx, id, UpdateAccountRequest{Name: "test updated"})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", account.Email)
	_, err = s.Update(ctx, "none", UpdateAccountRequest{Name: "test updated"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdateAccountRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateAccountRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	account, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", account.Email)
	assert.Equal(t, id, account.ID)

	// query
	accounts, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(accounts))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	account, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, account.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.Account
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Account, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.Account{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.Account, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, account entity.Account) error {
	if account.Email == "error" {
		return errCRUD
	}
	m.items = append(m.items, account)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, account entity.Account) error {
	if account.Email == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == account.ID {
			m.items[i] = account
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
