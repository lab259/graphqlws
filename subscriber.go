package graphqlws

// Topic represents a custom interface that represents a topic that will be used
// along with a PubSub system.
type Topic interface {
	// ID will return the structure ID of the topic on the technology used for
	// that purpose. For example, using a Redis PubSub system, this method would
	// return a string containing identifier of the channel.
	ID() interface{}
}

// StringTopic is a simple implementaiton of `Topic` for those PubSub systems
// that use simple strings as topics.
type StringTopic string

func (topic StringTopic) ID() interface{} {
	return topic
}

// Subscriber is reponsible for subscribing a Field to a topic.
type Subscriber interface {
	Subscribe(topic Topic) error
}
