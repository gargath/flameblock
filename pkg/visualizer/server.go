package visualizer

import (
	"fmt"
	"net/http"

	"github.com/gargath/flameblock/pkg/config"
	"github.com/go-redis/redis"
	"github.com/gobuffalo/packr/v2"
)

// Server handles the incoming webhooks
type Server struct {
	Config config.Configuration
	redis  *redis.Client
}

// Start will spin up the server and handle incoming webhooks
func (s *Server) Start() error {

	client, err := redisClientFromConfig(s.Config)
	if err != nil {
		return fmt.Errorf("Failed to connect to Redis: %v", err)
	}
	s.redis = client

	http.HandleFunc("/flamedata", s.flamedata)

	box := packr.New("static", "../../assets")

	buildHTTPHandlers(box)
	err = http.ListenAndServe(s.Config.BindAddr, nil)
	return fmt.Errorf("Error during ListenAndServe: %v", err)
}

func redisClientFromConfig(c config.Configuration) (*redis.Client, error) {
	var r *redis.Client
	if c.UseSentinel {
		r = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    c.SentinelMaster,
			SentinelAddrs: c.SentinelAddrs,
		})
	} else {
		r = redis.NewClient(&redis.Options{
			Addr:     c.RedisAddr,
			Password: "",
			DB:       0,
		})
	}
	pong, err := r.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis: Ping <> %v\n", pong)
	return r, nil
}
