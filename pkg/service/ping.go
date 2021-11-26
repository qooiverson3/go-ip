package service

import (
	"context"
	"fmt"
	"ipkeeper/pkg/model"
	"time"

	"github.com/go-ping/ping"
	"go.mongodb.org/mongo-driver/mongo"
)

type pingService struct {
	repo model.PingRepository
}

func NewPingService(r model.PingRepository) *pingService {
	return &pingService{
		repo: r,
	}
}

func (s *pingService) WriteResult(ctx context.Context, data model.IP) (*mongo.InsertOneResult, error) {
	return s.repo.WriteResult(ctx, data)
}

func (s *pingService) HealthCheck(input string) string {
	return input
}

func (s *pingService) GetAvailableIPs(r *model.IPResource) []model.IP {

	channel_data := make(chan []model.IP, 1)
	data := []model.IP{}
	channel_data <- data

	if 254*len(r.Networks) < r.Amount {
		return data
	}

	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < len(r.Networks); i++ {
		for n := 1; n < 254; n++ {
			go func(c context.Context, n string) {
				select {
				case <-c.Done():
					return
				default:
					d := <-channel_data
					if len(d) == r.Amount {
						channel_data <- d
						return
					}
					new_data, state := Ping(n, d)
					if state {
						channel_data <- new_data
						return
					}
					channel_data <- new_data
				}
			}(ctx, fmt.Sprintf("%s.%d", r.Networks[i], n))
		}
	}

	for {
		result := <-channel_data
		if len(result) == r.Amount {
			cancel()
			channel_data <- result
			break
		}
		channel_data <- result
	}

	return <-channel_data
}

func Ping(destination string, data []model.IP) ([]model.IP, bool) {

	pinger, err := ping.NewPinger(destination)
	if err != nil {
		return data, false
	}

	pinger.Count = 3
	pinger.Timeout = 200 * time.Millisecond
	if err := pinger.Run(); err != nil {
		return data, false
	}

	state := pinger.Statistics()
	if state.MaxRtt > 0 {
		ip := model.IP{}
		ip.IP = destination
		ip.RRT = state.MaxRtt
		ip.State = true

		data = append(data, ip)
		return data, true
	}
	return data, false
}
