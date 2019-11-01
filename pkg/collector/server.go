package collector

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gargath/flameblock/pkg/api"
)

// Server handles the incoming webhooks
type Server struct{}

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
	http.HandleFunc("/hook", s.hook)
	err := http.ListenAndServe(":8000", nil)
	return err
}
