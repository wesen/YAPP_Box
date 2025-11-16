package schemagen

import (
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
)

// DiscoverSchemaFiles walks modulesDir and returns all schema.yaml files.
func DiscoverSchemaFiles(modulesDir string) ([]string, error) {
	if modulesDir == "" {
		return nil, errors.New("modules directory is required")
	}

	var files []string
	if err := filepath.WalkDir(modulesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(path) == "schema.yaml" {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "walk modules dir %s", modulesDir)
	}

	sort.Strings(files)
	return files, nil
}

// LoadSchemaDocs validates and loads schema metadata.
func LoadSchemaDocs(paths []string) ([]*SchemaDoc, error) {
	var aggregated ValidationErrors
	var docs []*SchemaDoc
	for _, path := range paths {
		if err := ValidateSchemaFile(path); err != nil {
			if verrs, ok := err.(ValidationErrors); ok {
				aggregated = append(aggregated, verrs...)
				continue
			}
			return nil, errors.Wrapf(err, "validate schema %s", path)
		}
		doc, err := LoadSchemaDoc(path)
		if err != nil {
			return nil, errors.Wrapf(err, "load schema %s", path)
		}
		docs = append(docs, doc)
	}
	if len(aggregated) > 0 {
		return nil, aggregated
	}
	return docs, nil
}
