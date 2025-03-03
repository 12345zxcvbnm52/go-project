package proto

import (
	gin "github.com/gin-gonic/gin"
	httpserver "kenshop/goken/server/httpserver"
)

// service
// @BasePath /v1
// @Description This is a sample API
// @Host api.example.com
// @Title My API
// @Version 1.0.0
type MessagingHttpServer struct {
	Server *httpserver.Server
}

func RegisterMessagingHTTPServer(s *httpserver.Server) *MessagingHttpServer {
	ss := &MessagingHttpServer{
		Server: s,
	}
	return ss
}

// Update Message Summary | This function updates a message.
// @Accept application/json
// @Failure 500 {object} map[string]interface{}
// @Produce application/json
// @Router /messages/{message_id} [GET]
// @Success 200 {object} map[string]interface{}
// @Tags ken tag
// @Param message_id path string true "消息 ID"
// @Param message_name path string true "消息 Name"
func (s *MessagingHttpServer) UpdateMessage(c *gin.Context) {}

func (s *MessagingHttpServer) Execute() error {
	s.Server.Engine.GET("/messages/:message_id", s.UpdateMessage)
	return s.Server.Serve()
}
