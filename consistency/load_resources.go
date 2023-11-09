package consistency

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

func LoadResources(filesystem fs.FS) (resources []*unstructured.Unstructured, err error) {

	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	err = fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsRegular() {
			return nil
		}

		if !(filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".json") {
			return nil
		}

		f, err := filesystem.Open(path)
		if err != nil {
			return fmt.Errorf("could not open %s: %w", path, err)
		}

		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", path, err)
		}

		documents := bytes.Split(data, []byte("---"))

		for _, doc := range documents {
			if len(doc) == 0 {
				continue
			}

			obj := &unstructured.Unstructured{}
			_, _, err := decUnstructured.Decode(doc, nil, obj)
			if err != nil {
				return fmt.Errorf("error decoding %s: %w", path, err)
			}

			// Append the decoded object to the slice
			resources = append(resources, obj)
		}

		return nil
	})

	return resources, nil
}
