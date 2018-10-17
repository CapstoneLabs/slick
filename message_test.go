package slick

import (
	"testing"
	"time"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

var now = time.Now()

// Skeleton of a real message
var publicSlackMessage = slack.Msg{
	Text:      "This is a test.",
	User:      "U2147483697",
	Channel:   "C2147483705",
	Timestamp: string(now.Unix()),
}
var publicMessage = Message{
	Msg:        &publicSlackMessage,
	SubMessage: nil,
	MentionsMe: false,
	IsEdit:     false,
	FromMe:     false,
}
var privateSlackMessage = slack.Msg{
	Text:      "This is a test.",
	User:      "U2147483697",
	Channel:   "D2147483705",
	Timestamp: string(now.Unix()),
}
var privateMessage = Message{
	Msg:        &privateSlackMessage,
	SubMessage: nil,
	MentionsMe: false,
	IsEdit:     false,
	FromMe:     false,
}

func TestMessage(t *testing.T) {
	// IsPrivate
	assert.False(t, publicMessage.IsPrivate())
	assert.True(t, privateMessage.IsPrivate())
	// ContainsAnyCased
	assert.True(t, publicMessage.ContainsAnyCased([]string{"This", "not"}))
	assert.False(t, publicMessage.ContainsAnyCased([]string{"this", "not"}))
	// ContainsAny
	assert.True(t, publicMessage.ContainsAny([]string{"This", "not"}))
	assert.True(t, publicMessage.ContainsAny([]string{"this", "not"}))
	assert.False(t, publicMessage.ContainsAny([]string{"not", "here"}))
	// ContainsAll
	assert.True(t, publicMessage.ContainsAll([]string{"This", "is", "a", "test"}))
	assert.False(t, publicMessage.ContainsAll([]string{"This", "is not", "a", "test"}))
	// HasPrefix
	assert.True(t, publicMessage.HasPrefix("This"))
	assert.False(t, publicMessage.HasPrefix("this"))
}
