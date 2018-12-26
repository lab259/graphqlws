package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jamillosantos/handler"
	"github.com/lab259/graphql"
	"github.com/lab259/graphqlws"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}

type Message struct {
	Text string `json:"text"`
	User *User  `json:"user"`
}

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("Starting example server on :8085")

	users := make(map[string]*User)

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	messageType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Message",
		Fields: graphql.Fields{
			"text": &graphql.Field{
				Type: graphql.String,
			},
			"user": &graphql.Field{
				Type: userType,
			},
		},
	})

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
		"me": &graphql.Field{
			Type: userType,
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				return &User{
					Name: "unknkown",
				}, nil
			},
		},
	}}

	var subscriptionManager graphqlws.SubscriptionManager
	var broadcastMessage func(message *Message)

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "MutationRoot",
			Fields: graphql.Fields{
				"send": &graphql.Field{
					Args: graphql.FieldConfigArgument{
						"user": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"text": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Type: messageType,
					Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
						userName := p.Args["user"].(string)

						user, ok := users[userName]
						if !ok {
							return nil, fmt.Errorf("user '%s' not found", userName)
						}

						text, ok := p.Args["text"]
						if !ok {
							return nil, errors.New("you must pass the text")
						}
						m := &Message{
							Text: text.(string),
							User: user,
						}
						broadcastMessage(m)
						return m, nil
					},
				},
			},
		}),
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "SubscriptionRoot",
			Fields: graphql.Fields{
				"onJoin": &graphqlws.SubscriptionField{
					Field: graphql.Field{
						Type: userType,
						Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
							return p.Context.Value("onJoin.user").(*User), nil
						},
					},
					Subscribe: func(subscriber graphqlws.Subscriber) error {
						return subscriber.Subscribe(graphqlws.StringTopic("onJoin"))
					},
				},
				"onLeft": &graphqlws.SubscriptionField{
					Field: graphql.Field{
						Type: userType,
						Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
							return p.Context.Value("onLeft.user").(*User), nil
						},
					},
					Subscribe: func(subscriber graphqlws.Subscriber) error {
						return subscriber.Subscribe(graphqlws.StringTopic("onLeft"))
					},
				},
				"onMessage": &graphqlws.SubscriptionField{
					Field: graphql.Field{
						Type: messageType,
						Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
							message, ok := p.Context.Value("message").(*Message)
							if !ok {
								return nil, nil
							}
							return message, nil
						},
					},
					Subscribe: func(subscriber graphqlws.Subscriber) error {
						return subscriber.Subscribe(graphqlws.StringTopic("onMessage"))
					},
				},
			},
		}),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.WithField("err", err).Panic("GraphQL schema is invalid")
	}

	// Create subscription manager and GraphQL WS handler
	subscriptionManager = graphqlws.NewInMemorySubscriptionManager(&schema)

	// Sends this to all subscriptions
	broadcastJoin := func(user *User) {
		err := subscriptionManager.Publish(graphqlws.StringTopic("onJoin"), context.WithValue(context.Background(), "onJoin.user", user))
		if err != nil {
			log.Errorln(err.Error())
		}
	}

	broadcastLeft := func(user *User) {
		err := subscriptionManager.Publish(graphqlws.StringTopic("onLeft"), context.WithValue(context.Background(), "onLeft.user", user))
		if err != nil {
			log.Errorln(err.Error())
		}
	}

	broadcastMessage = func(message *Message) {
		err := subscriptionManager.Publish(graphqlws.StringTopic("onMessage"), context.WithValue(context.Background(), "message", message))
		if err != nil {
			log.Errorln(err.Error())
		}
	}

	connFactory := graphqlws.NewConnectionFactory(graphqlws.ConnectionConfig{
		Authenticate: func(token string) (interface{}, error) {
			if token == "Anonymous" {
				return nil, errors.New("forbidden")
			}
			user := &User{
				Name: token,
			}
			log.Infoln(token, "joined")
			users[user.Name] = user
			broadcastJoin(user)
			return user, nil
		},
		EventHandlers: graphqlws.ConnectionEventHandlers{
			Close: func(connection graphqlws.Connection) {
				u := connection.User().(*User)
				broadcastLeft(u)
				log.Infoln(u.Name, "left")
				delete(users, u.Name)
			},
		},
	}, subscriptionManager, *log.New())

	websocketHandler := graphqlws.NewHandler(graphqlws.HandlerConfig{
		ConnectionFactory:   connFactory,
		SubscriptionManager: subscriptionManager,
	})

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	// Serve the GraphQL WS endpoint
	mux := http.NewServeMux()
	mux.Handle("/graphql", h)
	mux.Handle("/subscriptions", websocketHandler)
	if err := http.ListenAndServe(":8085", cors.AllowAll().Handler(mux)); err != nil {
		log.WithField("err", err).Error("Failed to start server")
	}
}
