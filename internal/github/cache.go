package github

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v78/github"
	"github.com/stepbeta/talosctlenv/internal/utils"
)

// saveToCache saves release data to cache file
func saveToCache(allReleases []*github.RepositoryRelease) {
	releasesData, err := json.Marshal(CacheData{
		Timestamp: time.Now().UTC().String(),
		Releases:  allReleases,
	})
	if err != nil {
		fmt.Println("Failed to marshal release data to json:", err)
		return
	}
	cachePath, err := utils.GetCachePath()
	if err != nil {
		fmt.Println("Failed to retrieve cache path:", err)
		return
	}
	err = os.WriteFile(cachePath, []byte(releasesData), 0644)
	if err != nil {
		fmt.Println("Failed to save release data to cache:", err)
		return
	}
}

// readFromCache reads cached release data from file
func readFromCache() ([]*github.RepositoryRelease, error) {
	var err error
	cachePath, err := utils.GetCachePath()
	if err != nil {
		fmt.Println("Failed to retrieve cache path:", err)
		return nil, err
	}
	content, err := os.ReadFile(cachePath)
	if err != nil {
		fmt.Println("Failed to read cache data from file:", err)
		return nil, err
	}
	var cacheData CacheData
	err = json.Unmarshal(content, &cacheData)
	if err != nil {
		fmt.Println("Failed to unmarshal cache data:", err)
		return nil, err
	}
	// TODO check for cacheData.Timestamp age?
	return cacheData.Releases, nil
}

type CacheData struct {
	Timestamp string                      `json:"timestamp"`
	Releases  []*github.RepositoryRelease `json:"releases"`
}
