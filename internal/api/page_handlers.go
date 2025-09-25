package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 页面处理器

// indexPage 首页
func (s *Server) indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

// typesPage 类型管理页面
func (s *Server) typesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "types.html", gin.H{})
}

// thingsPage 事物管理页面
func (s *Server) thingsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "things.html", gin.H{})
}

// graphPage 关系图页面
func (s *Server) graphPage(c *gin.Context) {
	c.HTML(http.StatusOK, "graph.html", gin.H{})
}
