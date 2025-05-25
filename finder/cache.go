package finder

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func cacheFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, "file-finder", "files.json"), nil
}

func SaveToCache(files []string) error {
	path, err := cacheFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(files)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadFromCache() ([]string, error) {
	path, err := cacheFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var files []string
	if err := json.Unmarshal(data, &files); err != nil {
		return nil, err
	}
	return files, nil
}
