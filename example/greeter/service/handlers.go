package greeter

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/astranet/httpserve"
)

// This is a bonus part. Consider that your service already exposes some HTTP handler.
//
// :)

// Handler implements public API of the service.
type Handler interface {
	Check(c *httpserve.Context) httpserve.Response
}

type handler struct {
	Fingerprint string
}

// Check returns a simple message.
func (h *handler) Check(c *httpserve.Context) httpserve.Response {
	return httpserve.NewJSONResponse(http.StatusOK, map[string]string{
		"fingerprint": h.Fingerprint,
		"timestamp":   time.Now().Format(time.RFC3339),
		"status":      "ok",
	})
}

// NewHandler returns a new HTTP Handler of the service.
func NewHandler() Handler {
	return &handler{
		Fingerprint: getFingerprint(),
	}
}

// getFingerprint returns a random hexadecimal string, for testing responses with a roundrobin LB.
func getFingerprint() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}

// Add these two features, so meshRPC introspection could detect this handler:

var HandlerSpec Handler = &handler{}

// HTTPMethodsMap allows to specify what exact methods are allowed on this HTTP endpoint when exported
// to the meshRPC cluster.
func (h *handler) HTTPMethodsMap() map[string][]string {
	return map[string][]string{
		"Check": {
			"GET",
		},
	}
}
