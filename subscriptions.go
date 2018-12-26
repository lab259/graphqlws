package graphqlws

import (
	"context"
	"errors"
	"fmt"
	"github.com/lab259/graphql"
	"github.com/lab259/graphql/gqlerrors"
	"github.com/lab259/graphql/language/ast"
	"github.com/lab259/graphql/language/parser"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

// ErrorsFromGraphQLErrors convert from GraphQL errors to regular errors.
func ErrorsFromGraphQLErrors(errors []gqlerrors.FormattedError) []error {
	if len(errors) == 0 {
		return nil
	}

	out := make([]error, len(errors))
	for i := range errors {
		out[i] = errors[i]
	}
	return out
}

// SubscriptionSendDataFunc is a function that sends updated data
// for a specific subscription to the corresponding subscriber.
type SubscriptionSendDataFunc func(*DataMessagePayload)

type SubscriptionInterface interface {
	GetID() string
	GetQuery() string
	GetVariables() map[string]interface{}
	GetOperationName() string
	GetDocument() *ast.Document
	GetFields() []string
	GetConnection() Connection
	GetSendData() SubscriptionSendDataFunc
	SetFields(document []string)
	SetDocument(document *ast.Document)
}

// Subscription holds all information about a GraphQL subscription
// made by a client, including a function to send data back to the
// client when there are updates to the subscription query result.
type Subscription struct {
	ID            string
	Query         string
	Variables     map[string]interface{}
	OperationName string
	Document      *ast.Document
	Fields        []string
	Connection    Connection
	SendData      SubscriptionSendDataFunc
}

func (s *Subscription) GetID() string {
	return s.ID
}

func (s *Subscription) GetQuery() string {
	return s.Query
}

func (s *Subscription) GetVariables() map[string]interface{} {
	return s.Variables
}

func (s *Subscription) GetOperationName() string {
	return s.OperationName
}

func (s *Subscription) SetDocument(value *ast.Document) {
	s.Document = value
}

func (s *Subscription) GetDocument() *ast.Document {
	return s.Document
}

func (s *Subscription) SetFields(value []string) {
	s.Fields = value
}

func (s *Subscription) GetFields() []string {
	return s.Fields
}

func (s *Subscription) GetConnection() Connection {
	return s.Connection
}

func (s *Subscription) GetSendData() SubscriptionSendDataFunc {
	return s.SendData
}

// MatchesField returns true if the subscription is for data that
// belongs to the given field.
func (s *Subscription) MatchesField(field string) bool {
	if s.Document == nil || len(s.Fields) == 0 {
		return false
	}

	// The subscription matches the field if any of the queries have
	// the same name as the field
	for _, name := range s.Fields {
		if name == field {
			return true
		}
	}
	return false
}

// ConnectionSubscriptions defines a map of all subscriptions of
// a connection by their IDs.
type ConnectionSubscriptions map[string]SubscriptionInterface

// Subscriptions defines a map of connections to a map of
// subscription IDs to subscriptions.
type Subscriptions map[Connection]ConnectionSubscriptions

// SubscriptionManager provides a high-level interface to managing
// and accessing the subscriptions
// made by GraphQL WS clients.
type SubscriptionManager interface {
	// AddSubscription adds a new subscription to the manager.
	AddSubscription(Connection, SubscriptionInterface) []error

	// RemoveSubscription removes a subscription from the manager.
	RemoveSubscription(Connection, string)

	// RemoveSubscriptions removes all subscriptions of a client connection.
	RemoveSubscriptions(Connection)

	CreateSubscriptionSubscriber(subscription SubscriptionInterface) Subscriber

	Publish(topic Topic, ctx context.Context) error

	Subscribe(Subscriber) error
}

/**
 * The default implementation of the SubscriptionManager interface.
 */

type inMemorySubscriptionManager struct {
	schema  *graphql.Schema
	logger  *log.Entry
	topicsM sync.Mutex
	topics  map[Topic]map[string]SubscriptionInterface
}

func NewSubscriptionManagerWithLogger(schema *graphql.Schema, logger *log.Entry) SubscriptionManager {
	return newSubscriptionManager(schema, logger)
}

// NewSubscriptionManager creates a new subscription manager.
func NewInMemorySubscriptionManager(schema *graphql.Schema) SubscriptionManager {
	return newSubscriptionManager(schema, NewLogger("subscriptions"))
}

func newSubscriptionManager(schema *graphql.Schema, logger *log.Entry) SubscriptionManager {
	manager := new(inMemorySubscriptionManager)
	manager.topics = make(map[Topic]map[string]SubscriptionInterface)
	manager.logger = logger
	manager.schema = schema
	return manager
}

func (m *inMemorySubscriptionManager) Publish(topic Topic, ctx context.Context) error {
	subs, ok := m.topics[topic]
	if !ok {
		return nil
	}
	for _, sub := range subs {
		log.WithFields(log.Fields{
			"topic":          topic,
			"connID":         sub.GetConnection().ID(),
			"subscriptionID": sub.GetID(),
		}).Infoln("publishing")
		r := graphql.Execute(graphql.ExecuteParams{
			OperationName: sub.GetOperationName(),
			AST:           sub.GetDocument(),
			Schema:        *m.schema,
			Context:       ctx,
			Args:          sub.GetVariables(),
			Root:          nil,
		})
		sub.GetSendData()(&DataMessagePayload{
			Errors: ErrorsFromGraphQLErrors(r.Errors),
			Data:   r.Data,
		})
	}
	return nil
}

func (m *inMemorySubscriptionManager) Subscribe(sbsr Subscriber) error {
	m.topicsM.Lock()
	defer m.topicsM.Unlock()

	subscription := sbsr.Subscription()
	for _, topic := range sbsr.Topics() {
		subs, ok := m.topics[topic]
		if !ok {
			subs = make(map[string]SubscriptionInterface)
			m.topics[topic] = subs
		}
		_, ok = subs[subscription.GetID()]
		if !ok {
			subs[subscription.GetID()] = subscription
			log.WithFields(log.Fields{
				"connID":         subscription.GetConnection().ID(),
				"subscriptionID": subscription.GetID(),
				"topic":          topic,
			}).Infoln("subscribed")
		}
	}
	return nil
}

func (m *inMemorySubscriptionManager) AddSubscription(
	conn Connection,
	subscription SubscriptionInterface,
) []error {
	m.logger.WithFields(log.Fields{
		"conn":         conn.ID(),
		"subscription": subscription.GetID(),
	}).Info("Add subscription")

	if errors := validateSubscription(subscription); len(errors) > 0 {
		m.logger.WithField("errors", errors).Warn("Failed to add invalid subscription")
		return errors
	}

	// Parse the subscription query
	document, err := parser.Parse(parser.ParseParams{
		Source: subscription.GetQuery(),
	})
	if err != nil {
		m.logger.WithField("err", err).Warn("Failed to parse subscription query")
		return []error{err}
	}

	// Validate the query document
	validation := graphql.ValidateDocument(m.schema, document, nil)
	if !validation.IsValid {
		m.logger.WithFields(log.Fields{
			"errors": validation.Errors,
		}).Warn("Failed to validate subscription query")
		return ErrorsFromGraphQLErrors(validation.Errors)
	}

	// Remember the query document for later
	subscription.SetDocument(document)

	// Extract query names from the document (typically, there should only be one)
	subscription.SetFields(subscriptionFieldNamesFromDocument(document))

	subscriber := NewInMemorySubscriber(subscription)

	result := make([]error, 0)

	var fields graphql.Fields
	switch fs := m.schema.SubscriptionType().TypedConfig().Fields.(type) {
	case graphql.Fields:
		fields = fs
	case graphql.FieldsThunk:
		fields = fs()
	default:
		result = append(result, errors.New("fields type not supported"))
		return result
	}

	log.WithFields(log.Fields{
		"connID":         subscription.GetConnection().ID(),
		"subscriptionID": subscription.GetID(),
		"fields":         strings.Join(subscription.GetFields(), ", "),
	}).Infoln("subscribing")

	for _, fieldName := range subscription.GetFields() {
		field, ok := fields[fieldName]
		if !ok {
			panic(fmt.Sprintf("subscription %s not found", fieldName))
		}
		subscriptionField, ok := field.(*SubscriptionField)
		if !ok {
			panic(fmt.Sprintf("subscription %s is not a SubscriptionField", fieldName))
		}
		err = subscriptionField.Subscribe(subscriber)
		if err != nil {
			result = append(result, err)
			continue
		}
	}
	err = m.Subscribe(subscriber)
	if err != nil {
		result = append(result, err)
	}

	if len(result) > 0 {
		return result
	}
	return nil
}

func (m *inMemorySubscriptionManager) RemoveSubscription(
	conn Connection,
	subscriptionID string,
) {
	m.logger.WithFields(log.Fields{
		"conn":           conn.ID(),
		"subscriptionID": subscriptionID,
	}).Info("Remove subscription")

	for _, subs := range m.topics {
		delete(subs, subscriptionID)
	}
}

func (m *inMemorySubscriptionManager) RemoveSubscriptions(conn Connection) {
	m.logger.WithFields(log.Fields{
		"conn": conn.ID(),
	}).Info("Remove subscriptions")

	for _, subs := range m.topics {
		for key, subscription := range subs {
			if subscription.GetConnection() == conn {
				delete(subs, key)
			}
		}
	}
}

func (m *inMemorySubscriptionManager) CreateSubscriptionSubscriber(subscription SubscriptionInterface) Subscriber {
	return NewInMemorySubscriber(subscription)
}

func validateSubscription(s SubscriptionInterface) []error {
	errs := []error{}

	if s.GetID() == "" {
		errs = append(errs, errors.New("Subscription ID is empty"))
	}

	if s.GetConnection() == nil {
		errs = append(errs, errors.New("Subscription is not associated with a connection"))
	}

	if s.GetQuery() == "" {
		errs = append(errs, errors.New("Subscription query is empty"))
	}

	if s.GetSendData() == nil {
		errs = append(errs, errors.New("Subscription has no SendData function set"))
	}

	return errs
}
