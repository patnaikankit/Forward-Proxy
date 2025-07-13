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

func GetBlockedURLs() []string {
	Lock()
	urls := make([]string, 0, len(blocked))
	for url, isBlocked := range blocked {
		if isBlocked {
			urls = append(urls, url)
		}
	}
	return urls
}

func IsBlocked(url string) bool {
	Lock()
	return blocked[url]
}

// Helper functions for logging
func LogRequest(req *HTTPRequest) {
	Lock()
	var protocol string
	if req.IsHTTPS {
		protocol = "HTTPS"
	} else {
		protocol = "HTTP"
	}
	entry := "[" + time.Now().Format("15:04:05") + "] " + protocol + " --> " + req.Method + " " + req.URL + " " + req.Version

	// entry := time.Now().Format("15:04:05") + " - " + req.URL
	logs = append(logs, entry)
}

func GetLogs() []string {
	Lock()
	return logs
}
