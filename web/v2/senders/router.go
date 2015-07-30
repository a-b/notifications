package senders

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type collection interface {
	Set(conn models.ConnectionInterface, sender collections.Sender) (createdSender collections.Sender, err error)
	Get(conn models.ConnectionInterface, senderID, clientID string) (retrievedSender collections.Sender, err error)
}

type RouterConfig struct {
	RequestLogging    stack.Middleware
	Authenticator     stack.Middleware
	DatabaseAllocator stack.Middleware
	SendersCollection collections.SendersCollection
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()

	createStack := stack.NewStack(NewCreateHandler(config.SendersCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)
	getStack := stack.NewStack(NewGetHandler(config.SendersCollection)).Use(config.RequestLogging, config.Authenticator, config.DatabaseAllocator)

	router.Handle("/senders", createStack).Methods("POST").Name("POST /senders")
	router.Handle("/senders/{sender_id}", getStack).Methods("GET").Name("GET /senders/{sender_id}")

	return router
}
