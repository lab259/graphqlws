package graphqlws

import (
	"github.com/lab259/graphql"
	"github.com/lab259/graphql/language/parser"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gorilla/websocket"
)

// HandlerConfig stores the configuration of a GraphQL WebSocket handler.
type HandlerConfig struct {
	Schema              *graphql.Schema
	SubscriptionManager SubscriptionManager
	ConnectionFactory   ConnectionFactory
	Authenticate        AuthenticateFunc
	ConnectionEvents    ConnectionEventHandlers
}

// NewHandler creates a WebSocket handler for GraphQL WebSocket connections.
// This handler takes a SubscriptionManager and adds/removes subscriptions
// as they are started/stopped by the client.
func NewHandler(config HandlerConfig) http.Handler {
	// Create a WebSocket upgrader that requires clients to implement
	// the "graphql-ws" protocol
	var upgrader = websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: []string{"graphql-ws"},
	}

	logger := NewLogger("handler")

	// Create a map (used like a set) to manage client connections
	var connections = make(map[Connection]bool)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Establish a WebSocket connection
			var ws, err = upgrader.Upgrade(w, r, nil)

			// Bail out if the WebSocket connection could not be established
			if err != nil {
				logger.Warn("Failed to establish WebSocket connection", err)
				return
			}

			// Close the connection early if it doesn't implement the graphql-ws protocol
			if ws.Subprotocol() != "graphql-ws" {
				logger.Warn("Connection does not implement the GraphQL WS protocol")
				ws.Close()
				return
			}

			conn := config.ConnectionFactory.Create(ws, ConnectionEventHandlers{
				StartOperation: func(conn Connection, operationID string, payload *StartMessagePayload) []error {
					subscription := &Subscription{
						ID:            operationID,
						Query:         payload.Query,
						Variables:     payload.Variables,
						OperationName: payload.OperationName,
						Connection:    conn,
						SendData: func(data *DataMessagePayload) {
							conn.SendData(operationID, data)
						},
					}

					logger.WithFields(log.Fields{
						"conn":         conn.ID(),
						"subscription": subscription.GetID(),
					}).Info("Add subscription")

					if errors := ValidateSubscription(subscription); len(errors) > 0 {
						logger.WithField("errors", errors).Warn("Failed to add invalid subscription")
						return errors
					}

					// Parse the subscription query
					document, err := parser.Parse(parser.ParseParams{
						Source: subscription.GetQuery(),
					})
					if err != nil {
						logger.WithField("err", err).Warn("Failed to parse subscription query")
						return []error{err}
					}

					// Validate the query document
					validation := graphql.ValidateDocument(config.Schema, document, nil)
					if !validation.IsValid {
						logger.WithFields(log.Fields{
							"errors": validation.Errors,
						}).Warn("Failed to validate subscription query")
						return ErrorsFromGraphQLErrors(validation.Errors)
					}

					// Remember the query document for later
					subscription.SetDocument(document)

					// Extract query names from the document (typically, there should only be one)
					subscription.SetFields(SubscriptionFieldNamesFromDocument(document))

					return config.SubscriptionManager.AddSubscription(conn, subscription)
				},
				StopOperation: func(conn Connection, operationID string) {
					config.SubscriptionManager.RemoveSubscription(conn, operationID)
				},
				Close: func(conn Connection) {
					config.SubscriptionManager.RemoveSubscriptions(conn)
					delete(connections, conn)
				},
			})

			connections[conn] = true
		},
	)
}
