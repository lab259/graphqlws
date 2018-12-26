package graphqlws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// HandlerConfig stores the configuration of a GraphQL WebSocket handler.
type HandlerConfig struct {
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
					return config.SubscriptionManager.AddSubscription(conn, &Subscription{
						ID:            operationID,
						Query:         payload.Query,
						Variables:     payload.Variables,
						OperationName: payload.OperationName,
						Connection:    conn,
						SendData: func(data *DataMessagePayload) {
							conn.SendData(operationID, data)
						},
					})
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
