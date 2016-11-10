// Taken from https://github.com/gin-gonic/contrib/blob/master/rest/rest.go
// with some enhancements

package rest

import (
	"strings"

	gin "gopkg.in/gin-gonic/gin.v1"
)

// All of the methods are the same type as HandlerFunc
// if you don't want to support any methods of CRUD, then don't implement it
type CreateSupported interface {
	CreateHandler(*gin.Context)
}
type ListSupported interface {
	ListHandler(*gin.Context)
}
type TakeSupported interface {
	TakeHandler(*gin.Context)
}
type UpdateSupported interface {
	UpdateHandler(*gin.Context)
}
type DeleteSupported interface {
	DeleteHandler(*gin.Context)
}

// It defines
//   POST: /path
//   GET:  /path
//   PUT:  /path/:path
//   POST: /path/:path
//
// And with hierarchy
//   POST: /path/:path/path2
//   GET:  /path/:path/path2
//   PUT:  /path/:path/path2/:path2
//   POST: /path/:path/path2/:path2
func CRUD(group *gin.RouterGroup, path string, resource interface{}) {
	pathNames := strings.Split(path, "/")
	last := pathNames[len(pathNames)-1]
	for i := 0; i < len(pathNames); i++ {
		if pathNames[i] == ":" {
			pathNames[i] = ":" + pathNames[i-1]
		}
	}
	path = "/" + strings.Join(pathNames, "/")

	if resource, ok := resource.(CreateSupported); ok {
		group.POST(path, resource.CreateHandler)
	}
	if resource, ok := resource.(ListSupported); ok {
		group.GET(path, resource.ListHandler)
	}
	if resource, ok := resource.(TakeSupported); ok {
		group.GET(path+"/:"+last, resource.TakeHandler)
	}
	if resource, ok := resource.(UpdateSupported); ok {
		group.PUT(path+"/:"+last, resource.UpdateHandler)
	}
	if resource, ok := resource.(DeleteSupported); ok {
		group.DELETE(path+"/:"+last, resource.DeleteHandler)
	}
}
