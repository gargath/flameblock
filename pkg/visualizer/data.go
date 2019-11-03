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

	_, fudge := req.URL.Query()["normalize"]

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
		var err error
		keys, cursor, err = s.redis.Scan(cursor, "*", 5).Result()
		if err != nil {
			fmt.Printf("Error from Scan: %v\n", err)
			// Log error and return 500
		}
		if len(keys) > 0 {

			values, err = s.redis.MGet(keys...).Result()
			if err != nil {
				fmt.Printf("Error from MGet: %v\n", err)
				// Log error and return 500
			}
			fmt.Printf("Processing %d keys\n", len(keys))
			for i, key := range keys {
				fmt.Printf("Merging key %s with value %v\n", key, values[i])
				flamedata = merge(flamedata, key, values[i], fudge)
			}
		}
		if cursor == 0 {
			break
		}
		//		flamedata = merge(flamedata, entries)

	}
	var total int64
	for _, c := range flamedata.Children {
		total += c.Value
	}
	flamedata.Value = total

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(flamedata)
}

//func merge(root *api.Node, data map[string]interface{}) *api.Node {
func merge(root *api.Node, key string, val interface{}, fudge bool) *api.Node {
	currentNode := root
	parsedval, _ := strconv.ParseInt(val.(string), 10, 32)
	if parsedval == 0 {
		fmt.Printf("-----ALERT!!! ZERO VALUE DETECTED-----\n")
	}
	stackLines := strings.Split(key, "*")
	for i, stackLine := range stackLines {
		if fudge && strings.Contains(stackLine, "node_modules") {
			continue
		}
		nodeForLine := findNode(currentNode, stackLine)
		if i == len(stackLines)-1 {
			nodeForLine.Value = parsedval
		}
		currentNode = nodeForLine
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
