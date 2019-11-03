package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gargath/flameblock/pkg/api"
)

// hook handles incoming webhooks
func (s *Server) hook(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// decode JSON body into data structure
	decoder := json.NewDecoder(req.Body)
	var hook api.NsolidHook
	err := decoder.Decode(&hook)
	if err != nil {
		fmt.Printf("WARN: Failed to parse hook payload: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ignore webhooks for events other than nsolid-process-blocked
	if hook.Event != "nsolid-process-blocked" {
		fmt.Printf("Ignoring unknown event type '%s'\n", hook.Event)
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Handling Hook (Blocked For: %d)\n", hook.BlockedFor)
	fmt.Printf("Stack:\n%v\n", hook.Stack)

	// spin off processing into goroutine. After this point, there's no reason to keep the client waiting
	go s.transformAndStore(hook.Stack, hook.BlockedFor)
	rw.WriteHeader(http.StatusNoContent)
}

// transformAndStore takes the parsed webhook payload, transforms the stack frames into Redis keys and stores them in Redis
func (s *Server) transformAndStore(stack string, blockedFor int64) {
	lines := transformStack(stack)
	var key strings.Builder
	// Walk over the reversed frames, store value in redis, then concatenate next frame
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

// transformStack transforms the stack frames into Redis keys
func transformStack(stack string) []string {
	// First create regexp to remove leading spaces and "at"
	pattern := regexp.MustCompile(`^\s*\bat\b\s`)
	// Then split the frames on newline
	lines := strings.Split(stack, "\n")
	out := make([]string, len(lines))
	// Finally reverse the slice, eliminating the regexp matches along the way
	for i := len(lines) - 1; i >= 0; i-- {
		line := pattern.ReplaceAllString(lines[i], "")
		if line != "" {
			out[len(out)-(i+1)] = line
		}
	}
	return out
}
