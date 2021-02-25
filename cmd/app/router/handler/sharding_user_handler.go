package handler

import (
	"net/http"
	. "server/pkg/db"

	"github.com/gin-gonic/gin"
)

type ShardingUserHandler struct {
	BaseHandler
}

// 创建User接口参数
type CreateShardingUserParam struct {
	Uid  int64  `form:"uid" json:"uid" xml:"uid" binding:"required"`
	Name string `form:"name" json:"name" xml:"name" binding:"required"`
}

// CreateUser action
func (h *ShardingUserHandler) CreateUser(c *gin.Context) {
	var params CreateShardingUserParam
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Bad Request", "err": err})
		return
	}
	ShardingDB.Exec(params.Uid, "INSERT INTO users(uid, name) VALUES(?, ?)", params.Uid, params.Uid)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": params})
}

// GetUser action
func (h *ShardingUserHandler) GetUser(c *gin.Context) {
	uid := c.Param("id")
	// u := CreateShardingUserParam{}
	// ShardingDB.Query(uid, "SELECT * FROM users WHERE UID = ?", uid).Scan(&u)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": uid})
}
