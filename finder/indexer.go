package finder

import (
	"os"
	"path/filepath"
	"strings"
)

func IndexFiles(root string, exclude []string) ([]string, error) {
	var files []string

	excluded := map[string]bool{}
	for _, e := range exclude {
		excluded[e] = true
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		for ex := range excluded {
			rel, err := filepath.Rel(ex, path)
			if err == nil && (rel == "." || !strings.HasPrefix(rel, "..")) {
				return filepath.SkipDir
			}
		}

		// Skip hidden files and directories
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
