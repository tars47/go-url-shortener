package urlshort

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type pathUrl struct {
	Path string
	Url  string
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//

func YAMLHandler(yamlFileName string, fallback http.Handler) http.HandlerFunc {
	file, err := os.Open(yamlFileName)
	if err != nil {
		return fallback.ServeHTTP
	}

	yml, err := io.ReadAll(file)
	if err != nil {
		return fallback.ServeHTTP
	}

	var ms []pathUrl
	if err := yaml.Unmarshal(yml, &ms); err != nil {
		fmt.Println(err)
		return fallback.ServeHTTP
	}

	pathsToUrls := make(map[string]string)
	for _, v := range ms {
		pathsToUrls[v.Path] = v.Url
	}

	return MapHandler(pathsToUrls, fallback)

}
