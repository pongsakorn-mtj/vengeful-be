package repository

import (
	"context"
	"time"
	"vengeful-be/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

type User struct {
	FirstName             string    `bson:"firstName"`
	LastName              string    `bson:"lastName"`
	PhoneNo               string    `bson:"phoneNo"`
	Email                 string    `bson:"email"`
	IsAcceptTnc           bool      `bson:"isAcceptTnc"`
	IsAcceptPrivacyPolicy bool      `bson:"isAcceptPrivacyPolicy"`
	CreatedAt             time.Time `bson:"createdAt"`
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) Create(ctx context.Context, req *models.RegisterRequest) error {
	user := User{
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		PhoneNo:               req.PhoneNo,
		Email:                 req.Email,
		IsAcceptTnc:           req.IsAcceptTnc,
		IsAcceptPrivacyPolicy: req.IsAcceptPrivacyPolicy,
		CreatedAt:             time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) GetAll(ctx context.Context, page, limit int64) ([]User, int64, error) {
	var users []User

	// Calculate skip for pagination
	skip := (page - 1) * limit

	// Set options for pagination
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}}) // Sort by createdAt desc

	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Find all users with pagination
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
