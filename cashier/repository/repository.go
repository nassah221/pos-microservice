package repository

import (
	"context"
	"fmt"

	"pos-microservices/cashier/model"
	db "pos-microservices/cashier/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type cashierRepository struct {
	db.Store
}

func NewRepository(s db.Store) CashierRepository {
	return &cashierRepository{s}
}

type CashierRepository interface {
	Create(ctx context.Context, cashier *model.Cashier) (string, error)
	GetByID(ctx context.Context, id string) (*model.Cashier, error)
	GetByEmail(ctx context.Context, email string) (*model.Cashier, error)
	GetAll(ctx context.Context) ([]*model.Cashier, error)
	Update(ctx context.Context, cashier *model.Cashier) error
	Delete(ctx context.Context, id string) error
}

func (r *cashierRepository) Create(ctx context.Context, cashier *model.Cashier) (string, error) {
	cashier.ID = primitive.NewObjectID()

	result, err := r.Collection().InsertOne(ctx, cashier)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *cashierRepository) GetByID(ctx context.Context, id string) (*model.Cashier, error) {
	var cashier model.Cashier

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id to object id: %v", err)
	}

	result := r.Collection().FindOne(ctx, bson.D{{"_id", objID}})
	if result.Err() != nil {
		return nil, result.Err()
	}

	if err := result.Decode(&cashier); err != nil {
		return nil, fmt.Errorf("failed to decode cashier: %v", err)
	}

	return &cashier, nil
}

func (r *cashierRepository) GetByEmail(ctx context.Context, email string) (*model.Cashier, error) {
	var cashier model.Cashier

	result := r.Collection().FindOne(ctx, bson.D{{"email", email}})
	if result.Err() != nil {
		return nil, result.Err()
	}

	if err := result.Decode(&cashier); err != nil {
		return nil, fmt.Errorf("failed to decode cashier: %v", err)
	}

	return &cashier, nil
}

func (r *cashierRepository) GetAll(ctx context.Context) ([]*model.Cashier, error) {
	var cashiers model.Cashiers

	cur, err := r.Collection().Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var cashier model.Cashier

		if err := cur.Decode(&cashier); err != nil {
			return nil, err
		}
		cashiers = append(cashiers, &cashier)
	}

	return cashiers, nil
}

func (r *cashierRepository) Update(ctx context.Context, cashier *model.Cashier) error {
	result, err := r.Collection().ReplaceOne(ctx, bson.M{"_id": cashier.ID}, cashier)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *cashierRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert id to object id: %v", err)
	}

	result := r.Collection().FindOneAndDelete(ctx, bson.M{"_id": objID})
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
