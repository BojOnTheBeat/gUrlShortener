package main

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

//PathToURL just stores what a parsed yaml should look like
type PathToURL struct {
	Path string
	URL  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path

		newURL, ok := pathsToUrls[reqPath] // ok is true if newURL exists in the map

		if ok {

			fmt.Println("redirecting to " + newURL)

			http.Redirect(w, r, newURL, http.StatusSeeOther)
		}

		fmt.Println("falling back ")

		fallback.ServeHTTP(w, r)

	})
}

func buildMap(paths []PathToURL) map[string]string {

	res := make(map[string]string)

	for _, pathStruct := range paths {
		res[pathStruct.Path] = pathStruct.URL
	}

	return res
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathsToURLS := []PathToURL{}

	err := yaml.Unmarshal(yml, &pathsToURLS)

	if err != nil {
		return nil, err
	}

	yamlHandler := MapHandler(buildMap(pathsToURLS), fallback)

	return yamlHandler, nil
}
