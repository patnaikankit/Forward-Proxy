package utils

import (
	"sync"
	"time"
)

var (
	cache   = make(map[string][]byte)
	blocked = make(map[string]bool)
	logs    []string
	lock    sync.Mutex
)

func Lock() {
	lock.Lock()
	defer lock.Unlock()
}

// Helper functions for caching
func GetCache(url string) []byte {
	Lock()
	return cache[url]
}

func SetCache(url string, data []byte) {
	Lock()
	cache[url] = data
}

// Helper functions for blocking
func BlockURL(url string) {
	Lock()
	blocked[url] = true
}

func UnblockURL(url string) {
	Lock()
	blocked[url] = false
}

func IsBlocked(url string) bool {
	Lock()
	return blocked[url]
}

// Helper functions for logging
func LogRequest(req *HTTPRequest) {
	Lock()
	entry := time.Now().Format("15:04:05") + " - " + req.URL
	logs = append(logs, entry)
}

func GetLogs() []string {
	Lock()
	return logs
}
