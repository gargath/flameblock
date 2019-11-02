package visualizer

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

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
		if strings.HasSuffix(url, "index.html") {
			url = strings.TrimSuffix(url, "index.html")
		}
		url = "/" + url

		http.HandleFunc(url, func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-type", extensionToContentType[filepath.Ext(path)])
			bytes, err := box.Find(path)
			if err != nil {
				fmt.Printf("Error finding path %s in box: %v", path, err)
			}
			_, _ = writer.Write(bytes)
		})
	}
}
