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

// behaviorsPage 行为管理页面
func (s *Server) behaviorsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "behaviors.html", gin.H{})
}

// relationshipsPage 关系管理页面
func (s *Server) relationshipsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "relationships.html", gin.H{})
}

// actorsPage Actor管理页面
func (s *Server) actorsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "actors.html", gin.H{})
}

// testThemePage 主题切换测试页面
func (s *Server) testThemePage(c *gin.Context) {
	c.HTML(http.StatusOK, "test-theme.html", gin.H{})
}
