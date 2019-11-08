package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	RegionHeaderIdentifier = "X-Forwarded-Region"
	UserHeaderIdentifier   = "X-Forwarded-User"
)

// ResourceNotFound represents there is no request data.
type ResourceNotFound struct {
	Err error
}

func (r ResourceNotFound) Error() string {
	return r.Err.Error()
}

// BadDBError represents there is a db error, whatever the redis, etcd, mysql etc.
type BadDBResponse struct {
	Name string
	Err  error
}

func (b BadDBResponse) Error() string {
	return "DB: " + b.Name + " : " + b.Err.Error()
}

// InvalidRequest represents request args are not valid.
type InvalidRequest struct {
	Err error
}

func (i InvalidRequest) Error() string {
	return i.Err.Error()
}

func GenHttpCode(err error) int {

	switch err.(type) {
	case InvalidRequest:
		return http.StatusBadRequest
	case BadDBResponse:
		return http.StatusBadGateway
	case ResourceNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func GetRegionByHeader(h http.Header) string {
	if region, ok := h[RegionHeaderIdentifier]; ok {
		return region[0]
	}
	return ""
}

func GetUserByHeader(h http.Header) string {
	if region, ok := h[UserHeaderIdentifier]; ok {
		return region[0]
	}
	return ""
}

func StdError(msg string) map[string]interface{} {
	return gin.H{"message": msg}
}

func StdErrorf(format string, a ...interface{}) map[string]interface{} {
	return StdError(fmt.Sprintf(format, a...))
}
