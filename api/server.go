package api

import (
	db "github.com/atselvan/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type (
	// Server serves HTTP requests for simple bank.
	Server struct {
		store  *db.Store
		router *gin.Engine
	}
)

// NewServer creates a new HTTP server and set up the router.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
