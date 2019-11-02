package visualizer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gargath/flameblock/pkg/api"
)

func (s *Server) flamedata(rw http.ResponseWriter, req *http.Request) {

	var cursor uint64

	flamedata := &api.Node{
		Name:     "root",
		Value:    0,
		Children: []*api.Node{},
	}

	fmt.Printf("Processing new request\n")

	for {
		var keys []string
		var values []interface{}
		var entries = make(map[string]interface{})
		var err error
		keys, cursor, err = s.redis.Scan(cursor, "*", 5).Result()
		if cursor == 0 {
			break
		}
		if err != nil {
			fmt.Printf("Error from Scan: %v\n", err)
			// Log error and return 500
		}
		values, err = s.redis.MGet(keys...).Result()
		if err != nil {
			fmt.Printf("Error from MGet: %v\n", err)
			// Log error and return 500
		}

		for i, key := range keys {
			entries[key] = values[i]
		}
		flamedata = merge(flamedata, entries)

	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(flamedata)
}

func merge(root *api.Node, data map[string]interface{}) *api.Node {
	for k, v := range data {
		currentNode := root
		parsedval, _ := strconv.ParseInt(v.(string), 10, 32)
		stackLines := strings.Split(k, "*")
		for i, stackLine := range stackLines {
			nodeForLine := findNode(currentNode, stackLine)
			if i == len(stackLines)-1 {
				nodeForLine.Value = parsedval
			}
			currentNode = nodeForLine
		}
	}
	return root
}

func findNode(root *api.Node, name string) *api.Node {
	for _, c := range root.Children {
		if c.Name == name {
			return c
		}
	}
	newChild := &api.Node{
		Name:     name,
		Children: []*api.Node{},
	}
	root.Children = append(root.Children, newChild)
	return newChild
}
