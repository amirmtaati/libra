package handlers

import (
	"fmt"
	"github.com/amirmtaati/libra/internal/app"
	"github.com/amirmtaati/libra/pkg/scanner"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Server struct {
	app    *app.App
	router *gin.Engine
}

func NewServer(app *app.App) *Server {
	s := &Server{
		app:    app,
		router: gin.Default(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	{
		api := s.router.Group("/api")
		api.GET("/books", s.getBooks)
        api.POST("/scan", s.scanLibrary)      
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) getBooks(c *gin.Context) {
	//	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	books, err := s.app.BookService.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"books":  books,
		"total":  len(books),
		"limit":  limit,
		"offset": offset,
	})
}

func (s *Server) scanLibrary(c *gin.Context) {
	var req struct {
		Path string `json:"path"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if req.Path != "" {
		s.app.Scanner = scanner.NewScanner(req.Path)
	}

	go func() {
		if err := s.app.ScanLibrary(); err != nil {
			fmt.Printf("Background scan error: %v\n", err)
		}
	}()

	c.JSON(202, gin.H{"message": "Scan started"})
}
