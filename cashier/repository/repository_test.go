package repository

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"pos-microservices/cashier/model"
	db "pos-microservices/cashier/mongo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var r repoWrapper

type repoWrapper struct {
	CashierRepository
	Drop func(context.Context) error
}

func init() {
	cfg, err := db.NewConfig("../.env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	s, err := db.NewStore(context.Background(), cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	r = repoWrapper{NewRepository(s), s.Collection().Drop}
}

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(b)
}

func generateRandomEmail() string {
	return generateRandomString(10) + "@test.com"
}

func newCashier(t *testing.T) *model.Cashier {
	t.Helper()

	return &model.Cashier{
		Name:     fmt.Sprintf("%s-%s", t.Name(), generateRandomString(10)),
		Email:    generateRandomEmail(),
		Password: "foo",
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}
}

func TestCreateCashier(t *testing.T) {
	cashier := newCashier(t)

	id, err := r.Create(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)

	t.Cleanup(func() {
		r.Drop(context.Background())
	})
}

func TestGetCashierByID(t *testing.T) {
	cashier := newCashier(t)

	id, err := r.Create(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)

	found, err := r.GetByID(context.Background(), id)
	assert.NoError(t, err)

	assert.NotNil(t, found)
	assert.Equal(t, id, found.ID.Hex())

	// Negative test case
	id = primitive.NewObjectID().Hex()
	_, err = r.GetByID(context.Background(), id)
	assert.Error(t, err)

	t.Cleanup(func() {
		r.Drop(context.Background())
	})
}

func TestGetCashierByEmail(t *testing.T) {
	cashier := newCashier(t)

	email := "getbyemail@test.com"
	cashier.Email = email

	id, err := r.Create(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)

	found, err := r.GetByEmail(context.Background(), email)
	assert.NoError(t, err)

	assert.NotNil(t, found)
	assert.Equal(t, email, found.Email)

	t.Cleanup(func() {
		r.Drop(context.Background())
	})
}

func TestUpdateCashier(t *testing.T) {
	cashier := newCashier(t)

	id, err := r.Create(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)

	cashier.Name = "TEST_UPDATE"
	objId, _ := primitive.ObjectIDFromHex(id)

	cashier.ID = objId

	err = r.Update(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)

	found, _ := r.GetByID(context.Background(), id)
	assert.Equal(t, "TEST_UPDATE", found.Name)

	// Negative test case
	objId = primitive.NewObjectID()
	cashier.ID = objId

	err = r.Update(context.Background(), cashier)
	assert.Error(t, err)

	t.Cleanup(func() {
		r.Drop(context.Background())
	})
}

func TestListCashiers(t *testing.T) {
	insertCashiers := []struct {
		cashier *model.Cashier
		err     error
	}{
		{newCashier(t), nil},
		{newCashier(t), nil},
		{newCashier(t), nil},
	}

	for _, tt := range insertCashiers {
		_, err := r.Create(context.Background(), tt.cashier)
		assert.NoError(t, err)
	}

	cashiers, err := r.GetAll(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, cashiers)
	assert.Equal(t, len(insertCashiers), len(cashiers))

	for _, cashier := range cashiers {
		assert.NotEmpty(t, cashier.ID.Hex())
	}

	t.Cleanup(func() {
		r.Drop(context.Background())
	})
}

func TestDeleteCashiers(t *testing.T) {
	cashier := newCashier(t)

	id, err := r.Create(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)

	err = r.Delete(context.Background(), id)
	assert.NoError(t, err)

	// Negative test case
	id = primitive.NewObjectID().Hex()
	err = r.Delete(context.Background(), id)
	assert.Error(t, err)

	t.Cleanup(func() {
		r.Drop(context.Background())
	})
}
