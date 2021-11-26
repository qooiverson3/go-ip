package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type IP struct {
	Vlan  int
	IP    string
	State bool
	RRT   time.Duration
}

type IPResource struct {
	Networks []string
	Amount   int
}

type PingRepository interface {
	WriteResult(ctx context.Context, data IP) (*mongo.InsertOneResult, error)
}

type PingService interface {
	WriteResult(ctx context.Context, data IP) (*mongo.InsertOneResult, error)
	HealthCheck(string) string
	GetAvailableIPs(r *IPResource) []IP
}
