package collector

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"

	"github.com/gargath/flameblock/pkg/api"
)

// Server handles the incoming webhooks
type Server struct {
	Config Configuration
	redis  *redis.Client
}

func (s *Server) hook(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var hook api.NsolidHook
	err := decoder.Decode(&hook)
	if err != nil {
		fmt.Printf("WARN: Failed to parse hook payload: %s\n", err)
	}
	if hook.Event != "nsolid-process-blocked" {
		fmt.Printf("Ignoring unknown event type '%s'\n", hook.Event)
		return
	}
	fmt.Printf("\n\n---Incoming Request:---\nEvent: %s\nBlocked For: %d\nStack:\n%s\n", hook.Event, hook.BlockedFor, hook.Stack)
}

// Start will spin up the server and handle incoming webhooks
func (s *Server) Start() error {

	client, err := redisClientFromConfig(s.Config)
	if err != nil {
		return fmt.Errorf("Failed to connect to Redis: %v", err)
	}
	s.redis = client

	http.HandleFunc("/hook", s.hook)
	err = http.ListenAndServe(":8000", nil)
	return fmt.Errorf("Error during ListenAndServe: %v", err)
}

func redisClientFromConfig(c Configuration) (*redis.Client, error) {
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
