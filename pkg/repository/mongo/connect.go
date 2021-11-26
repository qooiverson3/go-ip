package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInfo struct {
	UserName, Password, Host string
	Port                     int
}

type MongoConnect struct {
	Client *mongo.Client
	DB     *mongo.Database
	Err    error
}

func (m *MongoConnect) ConnectToMongo(i *MongoInfo) *MongoConnect {

	auth := options.Credential{
		Username: i.UserName,
		Password: i.Password,
	}

	m.Client, m.Err = mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", i.Host, i.Port)).SetConnectTimeout(5*time.Second).SetAuth(auth))
	if m.Err != nil {
		return m
	}

	if err := m.Client.Ping(context.TODO(), nil); err != nil {
		return m
	}

	return m
}
