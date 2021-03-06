package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
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
		if redir, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, redir, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	mp, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	return MapHandler(mp, fallback), nil
}

func parseYAML(yml []byte) (map[string]string, error) {
	mp := make(map[string]string)

	// Parse
	var data []pathURL
	err := yaml.Unmarshal(yml, &data)
	if err != nil {
		return nil, err
	}

	// Move data to map
	for _, item := range data {
		mp[item.Path] = item.URL
	}

	return mp, nil
}

type pathURL struct {
	Path string `yaml:"path,omitempty"`
	URL  string `yaml:"url,omitempty"`
}
