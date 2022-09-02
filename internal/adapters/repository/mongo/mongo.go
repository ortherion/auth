package mongo

import (
	"auth/internal/config"
	"auth/internal/domain/models"
	"auth/internal/utils"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func NewClientMongoDB(ctx context.Context, cfg *config.Configs) (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.Dsn)

	clientOptions.SetAuth(options.Credential{
		AuthSource: "auth",
		Username:   os.Getenv("MONGO_LOGIN"),
		Password:   os.Getenv("MONGO_PASS"),
	})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to MongoDB: %w", err)
	}

	return client.Database(cfg.MongoDB.Database).Collection(cfg.MongoDB.Collection), nil
}

type MongoRepo struct {
	Collection *mongo.Collection
}

func NewMongoRepo(collection *mongo.Collection) *MongoRepo {
	return &MongoRepo{Collection: collection}
}

func (r *MongoRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	query, err := r.Collection.Find(ctx, bson.D{})
	defer query.Close(ctx)

	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)
	for query.Next(ctx) {
		var user models.User
		err := query.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *MongoRepo) Get(ctx context.Context, id string) (*models.User, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	query := r.Collection.FindOne(ctx, bson.M{"_id": docId})

	var user models.User
	err = query.Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		err = models.ErrUserNotFound
	}

	return &user, err
}

func (r *MongoRepo) GetByName(ctx context.Context, login string) (*models.User, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	query := r.Collection.FindOne(ctx, bson.M{"login": login})

	var user models.User
	err := query.Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		err = models.ErrUserNotFound
	}

	return &user, err
}

func (r *MongoRepo) Insert(ctx context.Context, user *models.User) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	dataReq := bson.M{
		"login":         user.Login,
		"password":      user.Password,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"role":          user.Role,
		"creation_date": time.Now().Unix(),
	}

	_, err := r.Collection.InsertOne(ctx, dataReq)
	if err != nil {
		return err
	}

	return err
}

func (r *MongoRepo) Update(ctx context.Context, user *models.User) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	dataReq := bson.M{
		"$set": bson.M{
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	res := r.Collection.FindOneAndUpdate(ctx, bson.M{"login": user.Login}, dataReq)

	return res.Err()
}

func (r *MongoRepo) UpdatePassword(ctx context.Context, user *models.User) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	dataReq := bson.M{
		"$set": bson.M{
			"password": user.Password,
		},
	}

	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": user.ID}, dataReq)

	return err
}
