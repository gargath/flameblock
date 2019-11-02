package visualizer

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
)

var extensionToContentType = map[string]string{
	".html": "text/html",
	".js":   "application/javascript",
	".css":  "text/css",
}

func buildHTTPHandlers(box *packr.Box) {

	list := box.List()
	for index := range list {
		path := list[index]
		url := path
		url = "/" + url
		fmt.Printf("Registering handler for %s\n", url)
		http.HandleFunc(url, func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-type", extensionToContentType[filepath.Ext(path)])
			fmt.Printf("Finding: %s\n", request.URL.Path)
			bytes, err := box.Find(request.URL.Path)
			fmt.Printf("Found: %d bytes\n", len(bytes))
			if err != nil {
				fmt.Printf("Error finding path %s in box: %v", path, err)
				return
			}
			_, _ = writer.Write(bytes)
		})
	}
}
