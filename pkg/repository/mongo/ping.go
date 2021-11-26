package mongo

import (
	"context"
	"ipkeeper/pkg/model"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dbName         string = "ips"
	collectionName string = "loginsight"
)

type pingRepository struct {
	db     *mongo.Database
	client *mongo.Client
}

func NewPingRepository(db *mongo.Database, client *mongo.Client) *pingRepository {
	return &pingRepository{db, client}
}

func (r *pingRepository) WriteResult(ctx context.Context, data model.IP) (*mongo.InsertOneResult, error) {
	collection := r.client.Database(dbName).Collection(collectionName)
	return collection.InsertOne(context.TODO(), &data)
}
