package slick

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

var now = time.Now()

// Skeleton of a real message
var publicSlackMessage = slack.Msg{
	Text:      "This is a test.",
	User:      "U2147483698",
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
var publicSlackMessageWithMention = slack.Msg{
	Text:      "This is a test. <@U2147483697>",
	User:      "U2147483697",
	Channel:   "C2147483705",
	Timestamp: string(now.Unix()),
}
var publicMessageWithMention = Message{
	Msg:        &publicSlackMessageWithMention,
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

var mockedBot = &Bot{
	Myself: slack.UserDetails{
		ID: "U2147483697",
	},
}

func TestIfMessageIsPrivateOrNot(t *testing.T) {
	assert.False(t, publicMessage.IsPrivate())
	assert.True(t, privateMessage.IsPrivate())
}

func TestIfMessageContainsCaseSensitiveWord(t *testing.T) {
	assert.True(t, publicMessage.ContainsAnyCased([]string{"This", "not"}))
	assert.False(t, publicMessage.ContainsAnyCased([]string{"this", "not"}))
}

func TestIfMessageContainsNonCaseSensitiveWord(t *testing.T) {
	assert.True(t, publicMessage.ContainsAny([]string{"This", "not"}))
	assert.True(t, publicMessage.ContainsAny([]string{"this", "not"}))
	assert.False(t, publicMessage.ContainsAny([]string{"not", "here"}))
}

func TestIfMessageContainsAllNonCaseSensitiveWords(t *testing.T) {
	assert.True(t, publicMessage.ContainsAll([]string{"This", "is", "a", "test"}))
	assert.False(t, publicMessage.ContainsAll([]string{"This", "is not", "a", "test"}))
}

func TestIfMessageContainsAStringNonCaseSensitive(t *testing.T) {
	assert.True(t, publicMessage.Contains("This"))
	assert.False(t, publicMessage.Contains("not"))
}

func TestIfMessageHasPrefix(t *testing.T) {
	assert.True(t, publicMessage.HasPrefix("This"))
	assert.False(t, publicMessage.HasPrefix("this"))
}

func TestShouldConvertAMessageToString(t *testing.T) {
	str := publicMessage.String()
	assert.True(t, strings.Contains(str, "Msg:"))
	assert.True(t, strings.Contains(str, "User:"))
}

func TestShouldNotApplyAMentionToMe(t *testing.T) {
	assert.False(t, publicMessage.MentionsMe)
	publicMessage.applyMentionsMe(mockedBot)
	assert.False(t, publicMessage.MentionsMe)
}

func TestShouldApplyAMentionToMeDirectMessage(t *testing.T) {
	assert.False(t, privateMessage.MentionsMe)
	privateMessage.applyMentionsMe(mockedBot)
	assert.True(t, privateMessage.MentionsMe)
}

func TestShouldApplyAMentionToMePublicMessageWithMention(t *testing.T) {
	assert.False(t, publicMessageWithMention.MentionsMe)
	publicMessageWithMention.applyMentionsMe(mockedBot)
	assert.True(t, publicMessageWithMention.MentionsMe)
}

func TestShouldNotApplyFromMe(t *testing.T) {
	assert.False(t, publicMessage.FromMe)
	publicMessage.applyFromMe(mockedBot)
	assert.False(t, publicMessage.FromMe)
}

func TestShouldApplyFromMeDirectMessage(t *testing.T) {
	assert.False(t, privateMessage.FromMe)
	privateMessage.applyFromMe(mockedBot)
	assert.True(t, privateMessage.FromMe)
}

func TestShouldReturnTheSameStringUnformatted(t *testing.T) {
	s := "text"
	fs := Format(s)
	assert.Equal(t, s, fs)
}

func TestShouldReturnTheStringConcatenatedWithAnotherString(t *testing.T) {
	s1 := "text"
	s2 := "text2"
	fs := Format(s1, s2)
	assert.Equal(t, fs, fmt.Sprintf(s1, s2))
}

func TestShouldReturnTheStringConcatenatedWithAnotherContentType(t *testing.T) {
	s1 := "text"
	i := 2
	fs := Format(s1, i)
	assert.Equal(t, fs, fmt.Sprintf(s1, i))
}
