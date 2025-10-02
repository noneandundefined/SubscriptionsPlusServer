package middleware

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"subscriptionplus/server/pkg/httpx"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const (
	MaxBodySize        = 1 << 20
	MaxViolationsPerIP = 3
	IpBlockDuration    = 10 * time.Minute
)

var (
	BlockedIPs   = make(map[string]time.Time)
	IpViolations = make(map[string]int)
	Mutex        sync.Mutex
)

// --- сигнатуры SQL-инъекций ---
var sqlInjectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bDROP\s+TABLE\b`),
	regexp.MustCompile(`(?i)\bUNION\s+SELECT\b`),
	regexp.MustCompile(`(?i)\bSELECT\b.+\bFROM\b`),
	regexp.MustCompile(`(?i)\bINSERT\s+INTO\b`),
	regexp.MustCompile(`(?i)\bUPDATE\s+\w+\s+SET\b`),
	regexp.MustCompile(`(?i)\bDELETE\s+FROM\b`),

	// классика
	regexp.MustCompile(`(?i)'\s*OR\s*'1'\s*=\s*'1'`),
	regexp.MustCompile(`(?i)\bOR\s+1\s*=\s*1\b`),
	regexp.MustCompile(`(?i)(?:--|#)`),
	regexp.MustCompile(`(?i);\s*exec\s+xp_cmdshell`),
}

// --- сигнатуры XSS-атак ---
var xssPattern = regexp.MustCompile(`(?i)(<script.*?>.*?</script>|<.*?on\w+\s*=\s*['"]?.*?['"]?|javascript:)`)

// SecurityMiddleware middleware проверок данных от пользователей (XSS, SQLInjection, ...)
func SecurityMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)

			if isBlocked(ip) {
				fmt.Printf("[ACCESS] Blocked IP %s\n", ip)
				httpx.HttpResponse(w, r, http.StatusForbidden, "access denied")
				return
			}

			contentType := r.Header.Get("Content-Type")
			if strings.HasPrefix(contentType, "application/octet-stream") {
				next.ServeHTTP(w, r)
				return
			}

			// limit size
			r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

			var bodyContent string
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				data, err := io.ReadAll(r.Body)
				if err != nil {
					httpx.HttpResponse(w, r, http.StatusRequestEntityTooLarge, "the request text is too big")
					registerViolation(ip)
					return
				}

				bodyContent = string(data)
				r.Body = io.NopCloser(strings.NewReader(bodyContent))
			}

			// check query params
			for _, values := range r.URL.Query() {
				for _, val := range values {
					if isMalicious(val) {
						httpx.HttpResponse(w, r, http.StatusBadRequest, "malicious request parameter")
						registerViolation(ip)
						return
					}
				}
			}

			// checkb form
			r.ParseForm()
			for _, values := range r.Form {
				for _, val := range values {
					if isMalicious(val) {
						httpx.HttpResponse(w, r, http.StatusBadRequest, "malicious request parameter")
						registerViolation(ip)
						return
					}
				}
			}

			// check body
			if bodyContent != "" && isMalicious(bodyContent) {
				httpx.HttpResponse(w, r, http.StatusBadRequest, "malicious content in the body")
				registerViolation(ip)
				return
			}

			// check headers
			for _, values := range r.Header {
				for _, val := range values {
					if isMalicious(val) {
						httpx.HttpResponse(w, r, http.StatusBadRequest, "malicious content")
						registerViolation(ip)
						return
					}
				}
			}

			// check cookies
			for _, cookie := range r.Cookies() {
				if isMalicious(cookie.Value) {
					httpx.HttpResponse(w, r, http.StatusBadRequest, "malicious content")
					registerViolation(ip)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isMalicious(input string) bool {
	for _, re := range sqlInjectionPatterns {
		if re.MatchString(input) {
			return true
		}
	}

	return xssPattern.MatchString(input)
}

func getIP(r *http.Request) string {
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	}

	host, _, _ := net.SplitHostPort(ip)
	return host
}

func isBlocked(ip string) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	expiry, exists := BlockedIPs[ip]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(BlockedIPs, ip)
		delete(IpViolations, ip)

		return false
	}

	return true
}

func registerViolation(ip string) {
	Mutex.Lock()
	defer Mutex.Unlock()

	IpViolations[ip]++
	if IpViolations[ip] >= MaxViolationsPerIP {
		BlockedIPs[ip] = time.Now().Add(IpBlockDuration)
	}
}
