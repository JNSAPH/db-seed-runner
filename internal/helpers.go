package internal

import (
	"os"
	"path/filepath"
	"sort"
)

func GetSQLFiles(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var files []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		files = append(files, filepath.Join(directory, entry.Name()))
	}

	sort.Strings(files)

	return files, nil
}
