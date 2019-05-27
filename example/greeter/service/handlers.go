package greeter

import (
	"time"

	"github.com/gin-gonic/gin"
)

// This is a bonus part. Consider that your service already exposes some HTTP handler.
//
// :)

type Handler interface {
	Check(c *gin.Context)
}

type handler struct{}

func (h *handler) Check(c *gin.Context) {
	c.String(200, "All ok! %s", time.Now().Format(time.RFC3339))
}

func NewHandler() Handler {
	return &handler{}
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
