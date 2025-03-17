package proto

import (
	gin "github.com/gin-gonic/gin"
)

// service
// @BasePath /v1
// @Description This is a sample API
// @Host api.example.com
// @Title My API
// @Version 1.0.0

type ReqMessageForm struct {
	Age   int32   `json:"age" form:"age" uri:"age" header:"" binding:""`
	Name  string  `json:"name" form:"name" uri:"name" header:"" binding:""`
	Price float64 `json:"price" form:"price" uri:"price" header:"" binding:""`
}

// Update Message Summary | This function updates a message.
// @Accept application/json
// @Failure 500 {object} map[string]interface{}
// @Produce application/json
// @Router /messages/{message_id} [GET]
// @Success 200 {object} map[string]interface{}
// @Tags ken tag
// @Param age path uint32 true "消息 ID"
// @Param name path string true "消息 Name"
// @Param price path float64 true "jia eg"
func UpdateMessage(c *gin.Context) {
}
