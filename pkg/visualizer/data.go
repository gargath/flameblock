package visualizer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gargath/flameblock/pkg/api"
)

// flamedata handles requests for the flamegraph data structure
func (s *Server) flamedata(rw http.ResponseWriter, req *http.Request) {

	// Check whether we're supposed to normalize / fudge the data
	_, fudge := req.URL.Query()["normalize"]

	// Redis cursor used for Scan()
	var cursor uint64

	// Create root node
	flamedata := &api.Node{
		Name:     "root",
		Value:    0,
		Children: []*api.Node{},
	}

	fmt.Printf("Processing new request\n")

	// Scan Redis keys and retrieve values until we've seen all keys
	for {
		var keys []string
		var values []interface{}
		var err error
		keys, cursor, err = s.redis.Scan(cursor, "*", 5).Result()
		if err != nil {
			fmt.Printf("Error from Scan: %v\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Check whether there were any more keys for this cursor
		if len(keys) > 0 {
			// Get values for all keys retrieved in this iteration
			values, err = s.redis.MGet(keys...).Result()
			if err != nil {
				fmt.Printf("Error from MGet: %v\n", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Printf("Processing %d keys\n", len(keys))
			for i, key := range keys {
				fmt.Printf("Merging key %s with value %v\n", key, values[i])
				// Merge current keys/values into data structure
				flamedata = merge(flamedata, key, values[i], fudge)
			}
		}
		if cursor == 0 {
			// There are no more keys
			break
		}

	}
	// Finally calculate total value for root node.
	var total int64
	for _, c := range flamedata.Children {
		total += c.Value
	}
	flamedata.Value = total

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	// Return the data
	json.NewEncoder(rw).Encode(flamedata)
}

//func merge(root *api.Node, data map[string]interface{}) *api.Node {
func merge(root *api.Node, key string, val interface{}, fudge bool) *api.Node {
	currentNode := root
	parsedval, _ := strconv.ParseInt(val.(string), 10, 32)
	if parsedval == 0 {
		// This should not happen
		fmt.Printf("-----ALERT!!! ZERO VALUE DETECTED-----\n")
	}
	// Split the Redis key to reconstruct the stack from the original webhook
	stackLines := strings.Split(key, "*")
	for i, stackLine := range stackLines {
		if fudge && strings.Contains(stackLine, "node_modules") {
			// ignore node_modules frames if fudge is true
			continue
		}
		// Find (or create) the node for the current stack line
		nodeForLine := findNode(currentNode, stackLine)
		if i == len(stackLines)-1 {
			// This is a leaf node, so write its value
			nodeForLine.Value = parsedval
		}
		// keep walking the tree until all lines are done
		currentNode = nodeForLine
	}
	return root
}

// findNode finds or creates the tree node for the given name
func findNode(root *api.Node, name string) *api.Node {
	// Walk the tree until we find the node. If found, return it
	for _, c := range root.Children {
		if c.Name == name {
			return c
		}
	}
	// The node did not exist if we got here, so create it, link it into the tree and then return it
	newChild := &api.Node{
		Name:     name,
		Children: []*api.Node{},
	}
	root.Children = append(root.Children, newChild)
	return newChild
}
