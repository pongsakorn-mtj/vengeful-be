package database

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
	logger *logrus.Logger
}

func NewMongoDB(logger *logrus.Logger) (*MongoDB, error) {
	username := os.Getenv("MONGO_ROOT_USERNAME")
	password := os.Getenv("MONGO_ROOT_PASSWORD")
	dbName := os.Getenv("MONGODB_DATABASE")
	host := os.Getenv("MONGODB_HOST")

	// Debug logging for environment variables
	logger.WithFields(logrus.Fields{
		"username_set": username != "",
		"password_set": password != "",
		"dbName":       dbName,
		"host":         host,
	}).Debug("MongoDB environment variables")

	if username == "" || password == "" || dbName == "" || host == "" {
		return nil, fmt.Errorf("required MongoDB environment variables are not set (username=%v, dbName=%v, host=%v)",
			username != "", dbName, host)
	}

	// Construct the connection string for MongoDB Atlas
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority&authMechanism=SCRAM-SHA-1",
		url.QueryEscape(username),
		url.QueryEscape(password),
		host,
		dbName,
	)

	// Log safe connection string (without credentials)
	logSafeURI := getSafeMongoURI(mongoURI)
	logger.WithField("uri", logSafeURI).Info("Connecting to MongoDB")

	// Set up client options with credentials
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    "admin",
		Username:      username,
		Password:      password,
	}

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetAuth(credential)

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	// Ping the database
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"database": dbName,
		"host":     host,
	}).Info("Successfully connected to MongoDB")

	return &MongoDB{
		client: client,
		db:     client.Database(dbName),
		logger: logger,
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		m.logger.WithError(err).Error("Failed to disconnect from MongoDB")
		return err
	}

	m.logger.Info("Successfully disconnected from MongoDB")
	return nil
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

// getSafeMongoURI removes sensitive information from MongoDB URI for logging
func getSafeMongoURI(uri string) string {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "INVALID_URI"
	}

	if parsed.User != nil {
		username := parsed.User.Username()
		parsed.User = url.UserPassword(username, "REDACTED")
	}

	return parsed.String()
}
