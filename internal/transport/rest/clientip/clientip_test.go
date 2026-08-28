package clientip_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegoclair/go_boilerplate/internal/transport/rest/clientip"
	"github.com/stretchr/testify/assert"
)

func request(headers map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.7:54321"
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	return req
}

func TestExtractor_ReadsTheEdgeThatIsInFront(t *testing.T) {
	var seen string
	extract := clientip.Extractor(nil, func(header string) { seen = header })

	assert.Equal(t, "203.0.113.7", extract(request(map[string]string{"CF-Connecting-IP": "203.0.113.7"})))
	assert.Equal(t, "cf-connecting-ip", seen)
}

// The order is what a deployment relies on when more than one edge writes a
// header; picking by chance would make the same request answer differently.
func TestExtractor_PrefersTheMoreSpecificHeader(t *testing.T) {
	extract := clientip.Extractor(nil, nil)

	ip := extract(request(map[string]string{
		"X-Real-IP":        "198.51.100.9",
		"CF-Connecting-IP": "203.0.113.7",
	}))

	assert.Equal(t, "203.0.113.7", ip)
}

// A deployment behind another edge says so, rather than the package guessing.
func TestExtractor_HonoursTheConfiguredOrder(t *testing.T) {
	extract := clientip.Extractor([]string{"X-Real-IP"}, nil)

	ip := extract(request(map[string]string{
		"X-Real-IP":        "198.51.100.9",
		"CF-Connecting-IP": "203.0.113.7",
	}))

	assert.Equal(t, "198.51.100.9", ip)
}

// The socket names the proxy rather than the caller, but no caller can forge
// it — which is why it is the fallback and X-Forwarded-For is not.
func TestExtractor_FallsBackToTheSocket(t *testing.T) {
	var seen string
	extract := clientip.Extractor(nil, func(header string) { seen = header })

	assert.Equal(t, "10.0.0.7", extract(request(nil)))
	assert.Empty(t, seen, "nothing answered, and the counter has to be able to say so")
}

// A header present but empty is the same as absent: an edge that forwards the
// name without a value may not shadow the one behind it.
func TestExtractor_SkipsAnEmptyHeader(t *testing.T) {
	extract := clientip.Extractor(nil, nil)

	ip := extract(request(map[string]string{
		"CF-Connecting-IP": "   ",
		"X-Real-IP":        "198.51.100.9",
	}))

	assert.Equal(t, "198.51.100.9", ip)
}

// X-Forwarded-For is a chain the caller can prepend to, so it may never be
// consulted by default — that was the hole the framework's old default had.
func TestExtractor_NeverTrustsForwardedForByDefault(t *testing.T) {
	extract := clientip.Extractor(nil, nil)

	ip := extract(request(map[string]string{"X-Forwarded-For": "1.2.3.4"}))

	assert.Equal(t, "10.0.0.7", ip)
}
