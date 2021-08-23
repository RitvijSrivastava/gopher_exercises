package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if value, ok := pathsToUrls[path]; ok {
			http.Redirect(rw, r, value, http.StatusFound)
			return
		}
		fallback.ServeHTTP(rw, r)
	}
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
	parsedYAML, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYAML)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]T, error) {
	var parsedYAML []T

	err := yaml.Unmarshal(yml, &parsedYAML)
	if err != nil {
		return nil, err
	}
	return parsedYAML, nil
}

func buildMap(parsedYAML []T) map[string]string {
	yamlMap := make(map[string]string)

	for _, yml := range parsedYAML {
		yamlMap[yml.Path] = yml.URL
	}
	return yamlMap
}

type T struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}