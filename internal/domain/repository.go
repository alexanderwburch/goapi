package domain

import (
	"context"

	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// Repository encapsulates the logic to access domains from the data source.
type Repository interface {
	// Get returns the domain with the specified domain ID.
	Get(ctx context.Context, id string) (entity.Domain, error)
	// Count returns the number of domains.
	Count(ctx context.Context) (int, error)
	// Query returns the list of domains with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]entity.Domain, error)
	// Create saves a new domain in the storage.
	Create(ctx context.Context, domain entity.Domain) error
	// Update updates the domain with given ID in the storage.
	Update(ctx context.Context, domain entity.Domain) error
	// Delete removes the domain with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists domains in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new domain repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the domain with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Domain, error) {
	var domain entity.Domain
	err := r.db.With(ctx).Select().Model(id, &domain)
	return domain, err
}

// Create saves a new domain record in the database.
// It returns the ID of the newly inserted domain record.
func (r repository) Create(ctx context.Context, domain entity.Domain) error {
	return r.db.With(ctx).Model(&domain).Insert()
}

// Update saves the changes to an domain in the database.
func (r repository) Update(ctx context.Context, domain entity.Domain) error {
	return r.db.With(ctx).Model(&domain).Update()
}

// Delete deletes an domain with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	domain, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&domain).Delete()
}

// Count returns the number of the domain records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("domain").Row(&count)
	return count, err
}

// Query retrieves the domain records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Domain, error) {
	var domains []entity.Domain
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&domains)
	return domains, err
}
