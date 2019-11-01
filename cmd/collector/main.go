package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gargath/flameblock/pkg/collector"
	flag "github.com/spf13/pflag"
)

var (
	redisAddr           = flag.String("redis-addr", "", "Address of Redis to connect to")
	redisUseSentinel    = flag.Bool("redis-use-sentinel", true, "Use Sentinels to handle Redis connections")
	redisSentinelAddrs  = flag.StringSlice("redis-sentinels", []string{}, "Address of Redis to connect to")
	redisSentinelMaster = flag.String("redis-sentinel-master", "", "The Redis Master name to use")
	showVersion         = flag.Bool("version", false, "Show version and exit")
)

func main() {
	fmt.Printf("Flameblock Collector version %s starting...\n", VERSION)

	flag.Parse()

	if *showVersion {
		fmt.Printf("flameblock collector %s\n", VERSION)
		return
	}

	config, cerr := validateConfig()
	if cerr != nil {
		fmt.Printf("Error validating config: %v\n", cerr)
		os.Exit(1)
	}

	c := &collector.Server{
		Config: *config,
	}

	err := c.Start()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func validateConfig() (*collector.Configuration, error) {
	config := &collector.Configuration{}
	if *redisUseSentinel {
		if *redisAddr != "" {
			return config, fmt.Errorf("Cannot specify Redis address when redis-use-sentinel is true")
		}
		if len(*redisSentinelAddrs) == 0 || *redisSentinelMaster == "" {
			return config, fmt.Errorf("Both redis-sentinel-addrs and redis-sentinel-master are required when redis-use-sentinel is true")
		}
		fmt.Printf("Using Redis Sentinel config: %s - Master: %s\n", strings.Join(*redisSentinelAddrs, ","), *redisSentinelMaster)
	} else {
		if *redisAddr == "" {
			return config, fmt.Errorf("Redis address is required")
		}
		if len(*redisSentinelAddrs) > 0 || *redisSentinelMaster != "" {
			return config, fmt.Errorf("Cannot specify redis-sentinel-addrs orredis-sentinel-master when redis-use-sentinel is false")
		}
		fmt.Printf("Using Redis config: %s\n", *redisAddr)
	}
	config.UseSentinel = *redisUseSentinel
	config.SentinelAddrs = *redisSentinelAddrs
	config.SentinelMaster = *redisSentinelMaster
	config.RedisAddr = *redisAddr

	return config, nil
}
