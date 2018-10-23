---
title: home
---

Slick is a Go-based ChatOps bot for Slack.

## Features

* Plugin interface for chat messages
* Plugin-based HTTP handlers
* Simple API to reply to users
* Keeps an internal state of channels, users and their state.
* Listen for Reactions; take actions based on them (like buttons).
* Simple API to message users privately
* Simple API to update a previously sent message
* Simple API to delete bot messages after a given time duration.
* Easy plugin interface, listeners with criteria such as:
  * Messages directed to the bot only
  * Private or public messages
  * Listens for a duration or until a given `time.Time`
  * Selectively on a channel, or from a user
  * Expire listeners and unregister them dynamically
  * Supports listening for edits or not
  * Regexp match messages, or Contains checks
* Built-in KV store for data persistence (backed by BoltDB and JSON serialization)
* The bot has a mood (_happy_ and _hyper_) which changes randomly.. you can base some decisions on it, to spice up conversations.
* Supports listening for any Slack events (ChannelCreated, ChannelJoined, EmojiChanged, FileShared, GroupArchived, etc..)
* A PubSub system to facilitate inter-plugins (or chat-to-web) communications.
