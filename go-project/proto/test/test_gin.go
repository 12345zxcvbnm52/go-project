package proto

import (
	strings "strings"

	gin "github.com/gin-gonic/gin"
)

type MessagingHttpServer struct {
	PathsMap map[string]gin.HandlerFunc
}

func RegisterMessagingHTTPServer() *MessagingHttpServer {
	s := &MessagingHttpServer{
		PathsMap: make(map[string]gin.HandlerFunc),
	}
	return s
}

func (s *MessagingHttpServer) RegisterHandlerFunc(method string, path string, handler gin.HandlerFunc) {
	s.PathsMap[method+";"+path] = handler
}

// Update Message Summary | This function updates a message.
// @Accept application/json
// @Failure 500 Internal Server Error
// @Produce application/json
// @Router /api/v1/resource
// @Success 200 OK
// @Tags ken tag
func (s *MessagingHttpServer) UpdateMessage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Update Message",
	})
}

func (s *MessagingHttpServer) register() {
	s.RegisterHandlerFunc("PATCH", "/v1/messages/:message_id", s.UpdateMessage)
}

func (s *MessagingHttpServer) Execute(r gin.IRouter) {
	s.register()
	for k, v := range s.PathsMap {
		path := strings.Split(k, ";")
		path[0] = strings.ToUpper(path[0])
		switch path[0] {
		case "GET":
			r.GET(path[1], v)
		case "PUT":
			r.PUT(path[1], v)
		case "DELETE":
			r.DELETE(path[1], v)
		case "POST":
			r.POST(path[1], v)
		case "PATCH":
			r.PATCH(path[1], v)
		}
	}
}
