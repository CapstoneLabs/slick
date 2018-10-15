package slick

import (
	"encoding/json"
	"testing"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

var (
	simpleChannel = `{
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
	simpleGroup = `{
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
	simpleIM = `{
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
)

func unmarshalIM(j string) (*slack.IM, error) {
	im := &slack.IM{}
	if err := json.Unmarshal([]byte(j), &im); err != nil {
		return nil, err
	}
	return im, nil
}

func unmarshalChannel(j string) (*slack.Channel, error) {
	channel := &slack.Channel{}
	if err := json.Unmarshal([]byte(j), &channel); err != nil {
		return nil, err
	}
	return channel, nil
}

func unmarshalGroup(j string) (*slack.Group, error) {
	group := &slack.Group{}
	if err := json.Unmarshal([]byte(j), &group); err != nil {
		return nil, err
	}
	return group, nil
}

// These tests kind of suck. They should migrate to something like:
// https://github.com/nlopes/slack/blob/master/conversation_test.go

func TestChannelFromSlackGroup(t *testing.T) {
	want, err := unmarshalGroup(simpleGroup)
	if err != nil {
		log.WithError(err).Error("Error unmarshalling JSON.")
	}
	got := ChannelFromSlackGroup(*want)

	if got.ID != want.ID {
		log.Error("Value mismatch")
	}
	if got.IsOpen != want.IsOpen {
		log.Error("Value mismatch")
	}
	if got.LastRead != want.LastRead {
		log.Error("Value mismatch")
	}
	if got.Name != want.Name {
		log.Error("Value mismatch")
	}
	if got.Creator != want.Creator {
		log.Error("Value mismatch")
	}
	if len(got.Members) != len(want.Members) {
		log.Error("Value mismatch")
	}
	if got.IsArchived != want.IsArchived {
		log.Error("Value mismatch")
	}
	if got.Topic != want.Topic {
		log.Error("Value mismatch")
	}
	if got.Purpose != want.Purpose {
		log.Error("Value mismatch")
	}
	if got.IsGroup != want.IsGroup {
		log.Error("Value mismatch")
	}
}

func TestChannelFromSlackChannel(t *testing.T) {
	want, err := unmarshalChannel(simpleChannel)
	if err != nil {
		log.WithError(err).Error("Error unmarshalling JSON.")
	}
	got := ChannelFromSlackChannel(*want)

	if got.ID != want.ID {
		log.Error("Value mismatch")
	}
	if got.IsOpen != want.IsOpen {
		log.Error("Value mismatch")
	}
	if got.LastRead != want.LastRead {
		log.Error("Value mismatch")
	}
	if got.Name != want.Name {
		log.Error("Value mismatch")
	}
	if got.Creator != want.Creator {
		log.Error("Value mismatch")
	}
	if len(got.Members) != len(want.Members) {
		log.Error("Value mismatch")
	}
	if got.IsArchived != want.IsArchived {
		log.Error("Value mismatch")
	}
	if got.Topic != want.Topic {
		log.Error("Value mismatch")
	}
	if got.Purpose != want.Purpose {
		log.Error("Value mismatch")
	}
	if got.IsGroup != want.IsGroup {
		log.Error("Value mismatch")
	}
}

func TestChannelFromSlackIM(t *testing.T) {
	want, err := unmarshalIM(simpleIM)
	if err != nil {
		log.WithError(err).Error("Error unmarshalling JSON.")
	}
	got := ChannelFromSlackIM(*want)

	if got.ID != want.ID {
		log.Error("Value mismatch")
	}
	if got.IsOpen != want.IsOpen {
		log.Error("Value mismatch")
	}
	if got.LastRead != want.LastRead {
		log.Error("Value mismatch")
	}
	if got.Name != want.User {
		log.Error("Value mismatch")
	}
	if got.User != want.User {
		log.Error("Value mismatch")
	}
	if got.IsUserDeleted != want.IsUserDeleted {
		log.Error("Value mismatch")
	}
	if got.IsIM != want.IsIM {
		log.Error("Value mismatch")
	}
}
