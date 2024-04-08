package account

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/test"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "account")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	err = repo.Create(ctx, entity.Account{
		ID:         "test1",
		Email:      "account1",
		FirebaseId: "xyz",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, 1, count2-count)

	// get
	account, err := repo.Get(ctx, "test1")
	assert.Nil(t, err)
	assert.Equal(t, "account1", account.Email)
	_, err = repo.Get(ctx, "test0")
	assert.Equal(t, sql.ErrNoRows, err)

	// update
	err = repo.Update(ctx, entity.Account{
		ID:        "test1",
		Email:     "account1 updated",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	assert.Nil(t, err)
	account, _ = repo.Get(ctx, "test1")
	assert.Equal(t, "account1 updated", account.Email)

	// query
	accounts, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, len(accounts))

	// delete
	err = repo.Delete(ctx, "test1")
	assert.Nil(t, err)
	_, err = repo.Get(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
	err = repo.Delete(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
}
