package slick

import (
	"encoding/json"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

var simpleChannel = `{
		"id": "C024BE91L",
		"name": "fun",
		"is_channel": true,
		"created": 1360782804,
		"creator": "U024BE7LH",
		"is_archived": false,
		"is_general": false,
		"members": [
			"U024BE7LH"
		],
		"topic": {
			"value": "Fun times",
			"creator": "U024BE7LV",
			"last_set": 1369677212
		},
		"purpose": {
			"value": "This channel is for fun",
			"creator": "U024BE7LH",
			"last_set": 1360782804
		},
		"is_member": true,
		"last_read": "1401383885.000061",
		"unread_count": 0,
		"unread_count_display": 0
	}`

var simpleGroup = `{
		"id": "G024BE91L",
		"name": "secretplans",
		"is_group": true,
		"created": 1360782804,
		"creator": "U024BE7LH",
		"is_archived": false,
		"members": [
			"U024BE7LH"
		],
		"topic": {
			"value": "Secret plans on hold",
			"creator": "U024BE7LV",
			"last_set": 1369677212
		},
		"purpose": {
			"value": "Discuss secret plans that no-one else should know",
			"creator": "U024BE7LH",
			"last_set": 1360782804
		},
		"last_read": "1401383885.000061",
		"unread_count": 0,
		"unread_count_display": 0
	}`

var simpleIM = `{
		"id": "D024BFF1M",
		"is_im": true,
		"user": "U024BE7LH",
		"created": 1360782804,
		"is_user_deleted": false,
		"is_open": true,
		"last_read": "1401383885.000061",
		"unread_count": 0,
		"unread_count_display": 0
	}`

func unmarshalGroup(j string) (*slack.Group, error) {
	group := &slack.Group{}
	if err := json.Unmarshal([]byte(j), &group); err != nil {
		return nil, err
	}
	return group, nil
}

func assertChannelFromSlackGroup(t *testing.T, slackGroup slack.Group, channel Channel) {
	assert.Equal(t, slackGroup.ID, channel.ID)
	assert.Equal(t, slackGroup.IsOpen, channel.IsOpen)
	assert.Equal(t, slackGroup.LastRead, channel.LastRead)
	assert.Equal(t, slackGroup.Name, channel.Name)
	assert.Equal(t, slackGroup.Creator, channel.Creator)
	assert.Equal(t, slackGroup.Members, channel.Members)
	assert.Equal(t, slackGroup.IsArchived, channel.IsArchived)
	assert.Equal(t, slackGroup.Topic, channel.Topic)
	assert.Equal(t, slackGroup.Purpose, channel.Purpose)
	assert.Equal(t, slackGroup.IsGroup, channel.IsGroup)
}

func TestChannelFromSlackGroup(t *testing.T) {
	slackGroup, err := unmarshalGroup(simpleGroup)
	assert.Nil(t, err)
	channel := ChannelFromSlackGroup(*slackGroup)

	assertChannelFromSlackGroup(t, *slackGroup, channel)
}

func unmarshalChannel(j string) (*slack.Channel, error) {
	channel := &slack.Channel{}
	if err := json.Unmarshal([]byte(j), &channel); err != nil {
		return nil, err
	}
	return channel, nil
}

func assertChannelFromSlackChannel(t *testing.T, slackChannel slack.Channel, channel Channel) {
	assert.Equal(t, slackChannel.ID, channel.ID)
	assert.Equal(t, slackChannel.IsOpen, channel.IsOpen)
	assert.Equal(t, slackChannel.LastRead, channel.LastRead)
	assert.Equal(t, slackChannel.Name, channel.Name)
	assert.Equal(t, slackChannel.Creator, channel.Creator)
	assert.Equal(t, slackChannel.Members, channel.Members)
	assert.Equal(t, slackChannel.IsArchived, channel.IsArchived)
	assert.Equal(t, slackChannel.Topic, channel.Topic)
	assert.Equal(t, slackChannel.Purpose, channel.Purpose)
	assert.Equal(t, slackChannel.IsGroup, channel.IsGroup)
}

func TestChannelFromSlackChannel(t *testing.T) {
	slackChannel, err := unmarshalChannel(simpleChannel)
	assert.Nil(t, err)
	channel := ChannelFromSlackChannel(*slackChannel)

	assertChannelFromSlackChannel(t, *slackChannel, channel)
}

func unmarshalIM(j string) (*slack.IM, error) {
	im := &slack.IM{}
	if err := json.Unmarshal([]byte(j), &im); err != nil {
		return nil, err
	}
	return im, nil
}

func assertChannelFromSlackIM(t *testing.T, slackIM slack.IM, channel Channel) {
	assert.Equal(t, slackIM.ID, channel.ID)
	assert.Equal(t, slackIM.IsOpen, channel.IsOpen)
	assert.Equal(t, slackIM.User, channel.Name)
	assert.Equal(t, slackIM.User, channel.User)
	assert.Equal(t, slackIM.IsUserDeleted, channel.IsUserDeleted)
	assert.Equal(t, slackIM.IsIM, channel.IsIM)
}

func TestChannelFromSlackIM(t *testing.T) {
	slackIM, err := unmarshalIM(simpleIM)
	assert.Nil(t, err)
	channel := ChannelFromSlackIM(*slackIM)

	assertChannelFromSlackIM(t, *slackIM, channel)
}
