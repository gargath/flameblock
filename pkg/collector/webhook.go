package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gargath/flameblock/pkg/api"
)

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
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Printf("Handling Hook (Blocked For: %d)\n", hook.BlockedFor)
	go s.transformAndStore(hook.Stack, hook.BlockedFor)
	rw.WriteHeader(http.StatusNoContent)
}

func (s *Server) transformAndStore(stack string, blockedFor int64) {
	lines := transformStack(stack)
	var key strings.Builder
	for i, line := range lines {
		key.WriteString(line)
		_, err := s.redis.IncrBy(key.String(), blockedFor).Result()
		if err != nil {
			fmt.Printf("ERROR: Failed to increment Redis key ('%s'): %v", key.String(), err)
			return
		}
		if i < len(lines)-1 {
			key.WriteString("*")
		}
	}
}

func transformStack(stack string) []string {
	pattern := regexp.MustCompile(`^\s*\bat\b\s`)
	lines := strings.Split(stack, "\n")
	out := make([]string, len(lines))
	for i := len(lines) - 1; i >= 0; i-- {
		line := pattern.ReplaceAllString(lines[i], "")
		if line != "" {
			out[len(out)-(i+1)] = line
		}
	}
	return out
}
