package graphqlws

import "github.com/lab259/graphql"

type SubscriptionField struct {
	graphql.Field
	Subscribe func(subscriber Subscriber) error
}
