package ops

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ManifestProvider func() ([]interface{}, error)

// ManifestsFromString takes a single manifest that may or may not contain
// multiple resources and returns a ManifestProvider of the parsed resources
func ManifestsFromString(manifest string) (ManifestProvider, error) {
	manifests, err := parseManifest(strings.NewReader(manifest))
	if err != nil {
		return nil, err
	}

	return func() ([]interface{}, error) {
		return manifests, nil
	}, nil
}

func ManifestsFromPath(manifestPath string) (ManifestProvider, error) {
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("target path %s does not exist: %w", manifestPath, err)
	}
	filesInPath, err := ioutil.ReadDir(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in path %s: %w", manifestPath, err)
	}
	manifests := []interface{}{}
	for _, fileInfo := range filesInPath {
		parsed, err := openAndParseFile(filepath.Join(manifestPath, fileInfo.Name()))
		if err != nil {
			return nil, err
		}
		manifests = append(manifests, parsed...)
	}
	return func() ([]interface{}, error) {
		return manifests, nil
	}, nil
}

func openAndParseFile(name string) ([]interface{}, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseManifest(f)
}

func parseManifest(manifestBody io.Reader) ([]interface{}, error) {
	manifests := []interface{}{}
	decoder := yaml.NewDecoder(manifestBody)
	for {
		m := map[string]interface{}{}
		err := decoder.Decode(&m)
		if err == io.EOF {
			break
		}
		manifests = append(manifests, m)

	}
	return manifests, nil
}
