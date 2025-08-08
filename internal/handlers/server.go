package handlers

import (
	"fmt"
	"github.com/amirmtaati/libra/internal/app"
	"github.com/amirmtaati/libra/pkg/scanner"
	"github.com/gin-gonic/gin"
	"strconv"
	"path/filepath"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

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
	api := s.router.Group("/api")
	{
		api.GET("/books", s.getBooks)
		api.GET("/books/:id", s.getBook)
		api.DELETE("/books/:id", s.deleteBook)
		api.GET("/books/:id/download", s.downloadBook)

		api.GET("/shelves", s.getShelves)
		api.POST("/shelves", s.createShelf)
		api.GET("/shelves/:id", s.getShelf)
		api.DELETE("/shelves/:id", s.deleteShelf)
		api.POST("/shelves/:id/books/:bookId", s.addBookToShelf)
		api.DELETE("/shelves/:id/books/:bookId", s.removeBookFromShelf)
		api.GET("/shelves/:id/books", s.getShelfBooks)

		api.POST("/scan", s.scanLibrary)
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) getBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	books, err := s.app.BookService.GetAll()
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    books,
		Meta: &Meta{
			Total:  len(books),
			Limit:  limit,
			Offset: offset,
		},
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

func (s *Server) getBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid book ID",
		})
		return
	}

	book, err := s.app.BookService.GetByID(uint(id))
	if err != nil {
		c.JSON(404, APIResponse{
			Success: false,
			Error:   "Book not found",
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    book,
	})
}

func (s *Server) deleteBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid book ID",
		})
		return
	}

	err = s.app.BookService.Delete(uint(id))
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Book deleted"},
	})
}

func (s *Server) downloadBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid book ID",
		})
		return
	}

	book, err := s.app.BookService.GetByID(uint(id))
	if err != nil {
		c.JSON(404, APIResponse{
			Success: false,
			Error:   "Book not found",
		})
		return
	}

	filename := filepath.Base(book.FilePath)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.File(book.FilePath)
}

func (s *Server) getShelves(c *gin.Context) {
	shelves, err := s.app.ShelfService.GetAll()
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    shelves,
	})
}

func (s *Server) createShelf(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid request",
		})
		return
	}

	shelf, err := s.app.ShelfService.Create(req.Name)
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(201, APIResponse{
		Success: true,
		Data:    shelf,
	})
}

func (s *Server) getShelf(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid shelf ID",
		})
		return
	}

	shelf, err := s.app.ShelfService.GetByID(uint(id))
	if err != nil {
		c.JSON(404, APIResponse{
			Success: false,
			Error:   "Shelf not found",
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    shelf,
	})
}

func (s *Server) deleteShelf(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid shelf ID",
		})
		return
	}

	err = s.app.ShelfService.Delete(uint(id))
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Shelf deleted"},
	})
}

func (s *Server) addBookToShelf(c *gin.Context) {
	shelfID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid shelf ID",
		})
		return
	}

	bookID, err := strconv.ParseUint(c.Param("bookId"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid book ID",
		})
		return
	}

	err = s.app.ShelfService.AddBook(uint(shelfID), uint(bookID))
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Book added to shelf"},
	})
}

func (s *Server) removeBookFromShelf(c *gin.Context) {
	shelfID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid shelf ID",
		})
		return
	}

	bookID, err := strconv.ParseUint(c.Param("bookId"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid book ID",
		})
		return
	}

	err = s.app.ShelfService.RemoveBook(uint(shelfID), uint(bookID))
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Book removed from shelf"},
	})
}

func (s *Server) getShelfBooks(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, APIResponse{
			Success: false,
			Error:   "Invalid shelf ID",
		})
		return
	}

	books, err := s.app.ShelfService.GetAllBooks(uint(id))
	if err != nil {
		c.JSON(500, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Data:    books,
	})
}
