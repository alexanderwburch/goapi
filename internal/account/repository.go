package account

import (
	"context"

	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// Repository encapsulates the logic to access accounts from the data source.
type Repository interface {
	// Get returns the account with the specified account ID.
	Get(ctx context.Context, id string) (entity.Account, error)
	// Count returns the number of accounts.
	Count(ctx context.Context) (int, error)
	// Query returns the list of accounts with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]entity.Account, error)
	// Create saves a new account in the storage.
	Create(ctx context.Context, account entity.Account) error
	// Update updates the account with given ID in the storage.
	Update(ctx context.Context, account entity.Account) error
	// Delete removes the account with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists accounts in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new account repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the account with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Account, error) {
	var account entity.Account
	err := r.db.With(ctx).Select().Model(id, &account)
	return account, err
}

// Create saves a new account record in the database.
// It returns the ID of the newly inserted account record.
func (r repository) Create(ctx context.Context, account entity.Account) error {
	return r.db.With(ctx).Model(&account).Insert()
}

// Update saves the changes to an account in the database.
func (r repository) Update(ctx context.Context, account entity.Account) error {
	return r.db.With(ctx).Model(&account).Update()
}

// Delete deletes an account with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	account, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&account).Delete()
}

// Count returns the number of the account records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("account").Row(&count)
	return count, err
}

// Query retrieves the account records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Account, error) {
	var accounts []entity.Account
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&accounts)
	return accounts, err
}
