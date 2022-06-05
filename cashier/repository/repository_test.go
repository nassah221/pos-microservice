package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"pos-microservices/cashier/model"
	db "pos-microservices/cashier/mongo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var r CashierRepository

func init() {
	cfg, err := db.NewConfig("../.env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	s, err := db.NewStore(context.Background(), cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	r = NewRepository(s)
	s.Collection().Drop(context.Background())
}

func TestCreateCashier(t *testing.T) {
	cashier := &model.Cashier{
		Name:     "TEST",
		Email:    "create@test.com",
		Password: "whatever",
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}

	id, err := r.Create(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)
}

func TestGetCashierByID(t *testing.T) {
	cashier := &model.Cashier{
		Name:     "TEST",
		Email:    "getbyid@test.com",
		Password: "whatever",
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}

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
}

func TestGetCashierByEmail(t *testing.T) {
	cashier := &model.Cashier{
		Name:     "TEST",
		Password: "whatever",
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}

	email := "getbyemail@test.com"
	cashier.Email = email

	id, err := r.Create(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)

	found, err := r.GetByEmail(context.Background(), email)
	assert.NoError(t, err)

	assert.NotNil(t, found)
	assert.Equal(t, email, found.Email)
}

func TestUpdateCashier(t *testing.T) {
	cashier := &model.Cashier{
		Name:     "TEST",
		Email:    "update@test.com",
		Password: "whatever",
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}

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
}

func TestListCashiers(t *testing.T) {
	cashiers, err := r.GetAll(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, cashiers)

	for _, cashier := range cashiers {
		assert.NotEmpty(t, cashier.ID.Hex())
	}
}

func TestDeleteCashiers(t *testing.T) {
	cashier := &model.Cashier{
		Name:     "TEST",
		Email:    "delete@test.com",
		Password: "whatever",
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}
	id, err := r.Create(context.Background(), cashier)
	assert.NoError(t, err)

	assert.NotEmpty(t, id)

	err = r.Delete(context.Background(), id)
	assert.NoError(t, err)

	// Negative test case
	id = primitive.NewObjectID().Hex()
	err = r.Delete(context.Background(), id)
	assert.Error(t, err)
}
