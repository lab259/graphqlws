package graphqlws

import (
	"errors"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/lab259/graphql"
	log "github.com/sirupsen/logrus"
)

// Mock connection

type mockWebSocketConnection struct {
	user string
	id   string
}

func (c *mockWebSocketConnection) Conn() *websocket.Conn {
	return nil
}

func (c *mockWebSocketConnection) ID() string {
	return c.id
}

func (c *mockWebSocketConnection) User() interface{} {
	return c.user
}

func (c *mockWebSocketConnection) SendData(
	opID string,
	data *DataMessagePayload,
) {
	// Do nothing
}

func (c *mockWebSocketConnection) SendError(err error) {
	// Do nothing
}

// Tests

func TestMain(m *testing.M) {
	log.SetLevel(log.ErrorLevel)
	m.Run()
}

func TestSubscriptions_NewSubscriptionManagerCreatesInstance(t *testing.T) {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{})

	sm := NewInMemorySubscriptionManager(&schema)
	if sm == nil {
		t.Fatal("NewSubscriptionManager fails in creating a new instance")
	}
}

func TestSubscriptions_AddingInvalidSubscriptionsFails(t *testing.T) {
	/*
	TODO
	 */
	/*
   schema, _ := graphql.NewSchema(graphql.SchemaConfig{})
   sm := NewSubscriptionManager(&schema)

   conn := mockWebSocketConnection{
	   id: "1",
   }

   // Try adding a subscription with nothing set
   errors := sm.AddSubscription(&conn, &Subscription{})

   if len(errors) == 0 {
	   t.Error("AddSubscription does not fail when adding an empty subscription")
   }

   if len(sm.Subscriptions()) > 0 {
	   t.Fatal("AddSubscription unexpectedly adds empty subscriptions")
   }

   // Try adding a subscription with an invalid query
   errors = sm.AddSubscription(&conn, &Subscription{
	   Query: "<<<Fooo>>>",
   })

   if len(errors) == 0 {
	   t.Error("AddSubscription does not fail when adding an invalid subscription")
   }

   if len(sm.Subscriptions()) > 0 {
	   t.Fatal("AddSubscription unexpectedly adds invalid subscriptions")
   }

   // Try adding a subscription with a query that doesn't match the schema
   errors = sm.AddSubscription(&conn, &Subscription{
	   Query: "subscription { foo }",
   })

   if len(errors) == 0 {
	   t.Error("AddSubscription doesn't fail if the query doesn't match the schema")
   }

   if len(sm.Subscriptions()) > 0 {
	   t.Fatal("AddSubscription unexpectedly adds invalid subscriptions")
   }
   */
}

func TestSubscriptions_AddingValidSubscriptionsWorks(t *testing.T) {
	/*
	TODO
	 */
	/*
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"users": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		})})
	sm := NewSubscriptionManager(&schema)

	conn := mockWebSocketConnection{
		id: "1",
	}

	// Add a valid subscription
	sub1 := Subscription{
		ID:         "1",
		Connection: &conn,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	errors := sm.AddSubscription(&conn, &sub1)

	if len(errors) > 0 {
		t.Error(
			"AddSubscription fails adding valid subscriptions. Unexpected errors:",
			errors,
		)
	}

	if len(sm.Subscriptions()) != 1 ||
		len(sm.Subscriptions()[&conn]) != 1 ||
		sm.Subscriptions()[&conn]["1"] != &sub1 {
		t.Fatal("AddSubscription doesn't add valid subscriptions properly")
	}

	// Add another valid subscription
	sub2 := Subscription{
		ID:         "2",
		Connection: &conn,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	errors = sm.AddSubscription(&conn, &sub2)
	if len(errors) > 0 {
		t.Error(
			"AddSubscription fails adding valid subscriptions.",
			"Unexpected errors:", errors,
		)
	}
	if len(sm.Subscriptions()) != 1 ||
		len(sm.Subscriptions()[&conn]) != 2 ||
		sm.Subscriptions()[&conn]["1"] != &sub1 ||
		sm.Subscriptions()[&conn]["2"] != &sub2 {
		t.Fatal("AddSubscription doesn't add valid subscriptions properly")
	}
	*/
}

func TestSubscriptions_AddingSubscriptionsTwiceFails(t *testing.T) {
	/*
	TODO
	*/
	/*
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"users": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		})})
	sm := NewSubscriptionManager(&schema)

	conn := mockWebSocketConnection{
		id: "1",
	}

	// Add a valid subscription
	sub := Subscription{
		ID:         "1",
		Connection: &conn,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	sm.AddSubscription(&conn, &sub)

	// Try adding the subscription for a second time
	errors := sm.AddSubscription(&conn, &sub)

	if len(errors) == 0 {
		t.Error(
			"AddSubscription doesn't fail when adding subscriptions a second time.",
			"Unexpected errors:", errors,
		)
	}

	if len(sm.Subscriptions()) != 1 ||
		len(sm.Subscriptions()[&conn]) != 1 ||
		sm.Subscriptions()[&conn]["1"] != &sub {
		t.Fatal("AddSubscription unexpectedly adds subscriptions twice")
	}
	*/
}

func TestSubscriptions_RemovingSubscriptionsWorks(t *testing.T) {
	/*
	TODO
	 */
	/*
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"users": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		})})
	sm := NewSubscriptionManager(&schema)

	conn := mockWebSocketConnection{
		id: "1",
	}

	// Add two valid subscriptions
	sub1 := Subscription{
		ID:         "1",
		Connection: &conn,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	sm.AddSubscription(&conn, &sub1)
	sub2 := Subscription{
		ID:         "2",
		Connection: &conn,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	sm.AddSubscription(&conn, &sub2)

	// Remove the first subscription
	sm.RemoveSubscription(&conn, &sub1)

	// Verify that only one subscription is left
	if len(sm.Subscriptions()) != 1 || len(sm.Subscriptions()[&conn]) != 1 {
		t.Error("RemoveSubscription does not remove subscriptions")
	}

	// Remove the second subscription
	sm.RemoveSubscription(&conn, &sub2)

	// Verify that there are no subscriptions left
	if len(sm.Subscriptions()) != 0 {
		t.Error("RemoveSubscription does not remove subscriptions")
	}
	*/
}

func TestSubscriptions_RemovingSubscriptionsOfAConnectionWorks(t *testing.T) {
	/*
	TODO
	 */
	/*
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Subscription: graphql.NewObject(graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"users": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		})})
	sm := NewSubscriptionManager(&schema)

	conn1 := mockWebSocketConnection{id: "1"}
	conn2 := mockWebSocketConnection{id: "2"}

	// Add four valid subscriptions
	sub1 := Subscription{
		ID:         "1",
		Connection: &conn1,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	sm.AddSubscription(&conn1, &sub1)
	sub2 := Subscription{
		ID:         "2",
		Connection: &conn1,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	sm.AddSubscription(&conn1, &sub2)
	sub3 := Subscription{
		ID:         "1",
		Connection: &conn2,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	sm.AddSubscription(&conn2, &sub3)
	sub4 := Subscription{
		ID:         "2",
		Connection: &conn2,
		Query:      "subscription { users }",
		SendData: func(msg *DataMessagePayload) {
			// Do nothing
		},
	}
	sm.AddSubscription(&conn2, &sub4)

	// Remove subscriptions of the first connection
	sm.RemoveSubscriptions(&conn1)

	// Verify that only the subscriptions of the second connection remain
	if len(sm.Subscriptions()) != 1 ||
		len(sm.Subscriptions()[&conn2]) != 2 ||
		sm.Subscriptions()[&conn1] != nil {
		t.Error("RemoveSubscriptions doesn't remove subscriptions of connections")
	}

	// Remove subscriptions of the second connection
	sm.RemoveSubscriptions(&conn2)

	// Verify that there are no subscriptions left
	if len(sm.Subscriptions()) != 0 {
		t.Error("RemoveSubscriptions doesn't remove subscriptions of connections")
	}
	*/
}

var _ = Describe("Subscriptions", func() {
	Describe("SubscriptionManager", func() {
		var schema graphql.Schema

		BeforeEach(func() {
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

			schemaConfig := graphql.SchemaConfig{
				Query: graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
					"me": &graphql.Field{
						Type: userType,
						Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
							return nil, nil
						},
					},
				}}),
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
								text, ok := p.Args["text"]
								if !ok {
									return nil, errors.New("you must pass the text")
								}
								return map[string]interface{}{
									"text": text.(string),
									"user": map[string]interface{}{
										"name": "Snake Eyes",
									},
								}, nil
							},
						},
					},
				}),
				Subscription: graphql.NewObject(graphql.ObjectConfig{
					Name: "SubscriptionRoot",
					Fields: graphql.Fields{
						"onJoin": &SubscriptionField{
							Field: graphql.Field{
								Type: userType,
								Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
									return map[string]interface{}{
										"name": "Snake Eyes",
									}, nil
								},
							},
							Subscribe: func(subscriber Subscriber) error {
								return subscriber.Subscribe(StringTopic("onJoin"))
							},
						},
						"onLeft": &SubscriptionField{
							Field: graphql.Field{
								Type: userType,
								Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
									return map[string]interface{}{
										"name": "Snake Eyes",
									}, nil
								},
							},
							Subscribe: func(subscriber Subscriber) error {
								return subscriber.Subscribe(StringTopic("onLeft"))
							},
						},
						"onMessage": &SubscriptionField{
							Field: graphql.Field{
								Type: messageType,
								Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
									return map[string]interface{}{
										"text": "this is a text",
										"user": map[string]interface{}{
											"name": "Snake Eyes",
										},
									}, nil
								},
							},
							Subscribe: func(subscriber Subscriber) error {
								return subscriber.Subscribe(StringTopic("onMessage"))
							},
						},
					},
				}),
			}

			var err error
			schema, err = graphql.NewSchema(schemaConfig)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should subscribe to a topic", func() {
			manager := NewInMemorySubscriptionManager(&schema)
			conn := &mockWebSocketConnection{}
			subscription := &Subscription{
				ID:         "1",
				Connection: conn,
				Query: `
subscription T {
	onJoin {
		name
	}
}
`,
				SendData: func(*DataMessagePayload) {
					//
				},
			}
			err := manager.AddSubscription(conn, subscription)
			Expect(err).To(BeEmpty())
			m := manager.(*inMemorySubscriptionManager)
			Expect(m.topics).To(HaveLen(1))
			Expect(m.topics).To(HaveKey(StringTopic("onJoin")))
			Expect(m.topics[StringTopic("onJoin")]).To(HaveLen(1))
			Expect(m.topics[StringTopic("onJoin")]["1"]).To(Equal(subscription))
		})

		It("should subscribe to multiple topics", func() {
			manager := NewInMemorySubscriptionManager(&schema)
			conn := &mockWebSocketConnection{}
			subscription := &Subscription{
				ID:         "1",
				Connection: conn,
				Query: `
subscription T {
	onJoin {
		name
	}
	onLeft {
		name
	}
	onMessage {
		user {
			name
		}
		text
	}
}
`,
				SendData: func(*DataMessagePayload) {
					//
				},
			}
			err := manager.AddSubscription(conn, subscription)
			Expect(err).To(BeEmpty())
			m := manager.(*inMemorySubscriptionManager)
			Expect(m.topics).To(HaveLen(3))
			Expect(m.topics).To(HaveKey(StringTopic("onJoin")))
			Expect(m.topics[StringTopic("onJoin")]).To(HaveLen(1))
			Expect(m.topics[StringTopic("onJoin")]["1"]).To(Equal(subscription))
			Expect(m.topics).To(HaveKey(StringTopic("onLeft")))
			Expect(m.topics[StringTopic("onJoin")]).To(HaveLen(1))
			Expect(m.topics[StringTopic("onJoin")]["1"]).To(Equal(subscription))
			Expect(m.topics).To(HaveKey(StringTopic("onMessage")))
			Expect(m.topics[StringTopic("onJoin")]).To(HaveLen(1))
			Expect(m.topics[StringTopic("onJoin")]["1"]).To(Equal(subscription))
		})

		It("should subscribe to multiple topics", func() {
			manager := NewInMemorySubscriptionManager(&schema)
			conn := &mockWebSocketConnection{}
			subscription := &Subscription{
				ID:         "1",
				Connection: conn,
				Query: `
subscription T {
	onJoin {
		name
	}
	onLeft {
		name
	}
	onMessage {
		user {
			name
		}
		text
	}
}
`,
				SendData: func(*DataMessagePayload) {
					//
				},
			}
			err := manager.AddSubscription(conn, subscription)
			Expect(err).To(BeEmpty())
			m := manager.(*inMemorySubscriptionManager)
			Expect(m.topics).To(HaveLen(3))
			Expect(m.topics).To(HaveKey(StringTopic("onJoin")))
			Expect(m.topics[StringTopic("onJoin")]).To(HaveLen(1))
			Expect(m.topics[StringTopic("onJoin")]["1"]).To(Equal(subscription))
			Expect(m.topics).To(HaveKey(StringTopic("onLeft")))
			Expect(m.topics[StringTopic("onJoin")]).To(HaveLen(1))
			Expect(m.topics[StringTopic("onJoin")]["1"]).To(Equal(subscription))
			Expect(m.topics).To(HaveKey(StringTopic("onMessage")))
			Expect(m.topics[StringTopic("onJoin")]).To(HaveLen(1))
			Expect(m.topics[StringTopic("onJoin")]["1"]).To(Equal(subscription))
		})

		It("should subscribe multiple clients to multiple topics", func() {
			manager := NewInMemorySubscriptionManager(&schema)
			conn1 := &mockWebSocketConnection{}
			subscription1 := &Subscription{
				ID:         "1",
				Connection: conn1,
				Query: `
subscription T {
	onJoin {
		name
	}
	onLeft {
		name
	}
	onMessage {
		user {
			name
		}
		text
	}
}
`,
				SendData: func(*DataMessagePayload) {
					//
				},
			}
			err := manager.AddSubscription(conn1, subscription1)
			Expect(err).To(BeEmpty())
			conn2 := &mockWebSocketConnection{}
			subscription2 := &Subscription{
				ID:         "2",
				Connection: conn2,
				Query: `
subscription T {
	onLeft {
		name
	}
	onMessage {
		user {
			name
		}
		text
	}
}
`,
				SendData: func(*DataMessagePayload) {
					//
				},
			}
			err = manager.AddSubscription(conn1, subscription1)
			Expect(err).To(BeEmpty())
			err = manager.AddSubscription(conn2, subscription2)
			Expect(err).To(BeEmpty())
			m := manager.(*inMemorySubscriptionManager)
			Expect(m.topics).To(HaveLen(3))
			Expect(m.topics).To(HaveKey(StringTopic("onJoin")))
			Expect(m.topics[StringTopic("onLeft")]).To(HaveLen(2))
			Expect(m.topics[StringTopic("onLeft")]["1"]).To(Equal(subscription1))
			Expect(m.topics[StringTopic("onLeft")]["2"]).To(Equal(subscription2))
			Expect(m.topics).To(HaveKey(StringTopic("onLeft")))
			Expect(m.topics[StringTopic("onJoin")]).To(HaveLen(1))
			Expect(m.topics[StringTopic("onJoin")]["1"]).To(Equal(subscription1))
			Expect(m.topics).To(HaveKey(StringTopic("onMessage")))
			Expect(m.topics[StringTopic("onMessage")]).To(HaveLen(2))
			Expect(m.topics[StringTopic("onMessage")]["1"]).To(Equal(subscription1))
			Expect(m.topics[StringTopic("onMessage")]["2"]).To(Equal(subscription2))
		})
	})
})
