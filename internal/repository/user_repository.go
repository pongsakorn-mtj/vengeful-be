package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
	MarketingCode         string    `bson:"marketingCode"`
	MarketedBy            string    `bson:"marketedBy,omitempty"`
	CreatedAt             time.Time `bson:"createdAt"`
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) generateMarketingCode() (string, error) {
	bytes := make([]byte, 6) // Generate 6 random bytes
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}

func (r *UserRepository) GetUserByMarketingCode(ctx context.Context, code string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"marketingCode": code}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Create(ctx context.Context, req *models.RegisterRequest) (string, error) {
	// Generate unique marketing code
	var marketingCode string
	var err error
	for i := 0; i < 5; i++ { // Try up to 5 times to generate unique code
		marketingCode, err = r.generateMarketingCode()
		if err != nil {
			return "", fmt.Errorf("failed to generate marketing code: %w", err)
		}

		// Check if code already exists
		exists, err := r.collection.CountDocuments(ctx, bson.M{"marketingCode": marketingCode})
		if err != nil {
			return "", err
		}
		if exists == 0 {
			break
		}
		if i == 4 {
			return "", fmt.Errorf("failed to generate unique marketing code after 5 attempts")
		}
	}

	// Check marketer if marketing code is provided
	var marketedBy string
	if req.MarketingCode != "" {
		marketer, err := r.GetUserByMarketingCode(ctx, req.MarketingCode)
		if err != nil {
			return "", fmt.Errorf("failed to check marketing code: %w", err)
		}
		if marketer != nil {
			marketedBy = marketer.Email
		}
	}

	user := User{
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		PhoneNo:               req.PhoneNo,
		Email:                 req.Email,
		IsAcceptTnc:           req.IsAcceptTnc,
		IsAcceptPrivacyPolicy: req.IsAcceptPrivacyPolicy,
		MarketingCode:         marketingCode,
		MarketedBy:            marketedBy,
		CreatedAt:             time.Now(),
	}

	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return marketingCode, nil
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
