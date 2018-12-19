package graphqlws

import "github.com/graphql-go/graphql"

type SubscriptionField struct {
	graphql.Field
	Subscribe func(subscriber Subscriber) error
}
