package graphqlws

type inMemorySubscriber struct {
	subscription SubscriptionInterface
	topics       []Topic
}

const PUBLISHED_DATA = "PUBLISHED_DATA"

func NewInMemorySubscriber(subscription SubscriptionInterface) Subscriber {
	return &inMemorySubscriber{
		subscription: subscription,
	}
}

func (subscriber *inMemorySubscriber) Subscription() SubscriptionInterface {
	return subscriber.subscription
}

func (subscriber *inMemorySubscriber) Topics() []Topic {
	return subscriber.topics
}

func (subscriber *inMemorySubscriber) Subscribe(topic Topic) error {
	if subscriber.topics == nil {
		subscriber.topics = make([]Topic, 0, 3)
	}
	subscriber.topics = append(subscriber.topics, topic)
	return nil
}
