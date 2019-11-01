package collector

// Configuration contains the command line arguments passed to collector
type Configuration struct {
	UseSentinel    bool
	SentinelAddrs  []string
	SentinelMaster string
	RedisAddr      string
}
