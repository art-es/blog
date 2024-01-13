package api

import "github.com/gin-gonic/gin"

type EndpointHandler interface {
	Method() string
	Endpoint() string
	Handle(*gin.Context)
}

func BindEndpoint(r *gin.Engine, h EndpointHandler) {
	r.Handle(h.Method(), h.Endpoint(), h.Handle)
}

func BindEndpoints(r *gin.Engine, hh ...EndpointHandler) {
	for _, h := range hh {
		BindEndpoint(r, h)
	}
}
