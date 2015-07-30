package senders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/ryanmoran/stack"
)

type GetHandler struct {
	senders collection
}

func NewGetHandler(senders collection) GetHandler {
	return GetHandler{
		senders: senders,
	}
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	senderID := splitURL[len(splitURL)-1]

	if senderID == "" {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{ "error": "%s" }`, "missing sender id")
		return
	}

	clientID := context.Get("client_id")
	if clientID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{ "error": "%s" }`, "missing client id")
		return
	}

	database := context.Get("database").(models.DatabaseInterface)
	sender, err := h.senders.Get(database.Connection(), senderID, context.Get("client_id").(string))
	if err != nil {
		switch err.(type) {
		case collections.NotFoundError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{ "error": "sender not found" }`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "error": "%s" }`, err)
		}
		return
	}

	getResponse, _ := json.Marshal(map[string]string{
		"id":   sender.ID,
		"name": sender.Name,
	})

	w.Write(getResponse)
}
