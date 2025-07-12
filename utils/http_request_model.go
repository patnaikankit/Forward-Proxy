package utils

type HTTPRequest struct {
	Method  string
	URL     string
	Version string
	Host    string
	Port    string
	IsHTTPS bool
	Raw     []byte
}
