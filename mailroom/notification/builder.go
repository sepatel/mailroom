// Copyright 2024 SeatGeek, Inc.
//
// Licensed under the terms of the Apache-2.0 license. See LICENSE file in project root for terms.

package notification

import (
	"github.com/seatgeek/mailroom/mailroom/common"
	"github.com/seatgeek/mailroom/mailroom/event"
	"github.com/seatgeek/mailroom/mailroom/identifier"
	slack2 "github.com/seatgeek/mailroom/mailroom/notifier/slack"
	"github.com/slack-go/slack"
)

type builderOpts struct {
	eventID             event.ID
	eventType           event.Type
	recipients          identifier.Collection
	fallbackMessage     string
	messagePerTransport map[common.TransportKey]string
	slackOpts           []slack.MsgOption
}

type Builder struct {
	opts builderOpts
}

func NewBuilder(eventID event.ID, eventType event.Type) *Builder {
	return &Builder{
		opts: builderOpts{
			eventID:             eventID,
			eventType:           eventType,
			recipients:          identifier.NewCollection(),
			messagePerTransport: make(map[common.TransportKey]string),
		},
	}
}

// WithRecipient sets the recipient of the notification
// It's like WithRecipientIdentifiers but it accepts a single identifier collection
func (b *Builder) WithRecipient(identifiers identifier.Collection) *Builder {
	b.opts.recipients = identifiers
	return b
}

// WithRecipientIdentifiers sets the recipient of the notification
// It's like WithRecipient but it accepts multiple identifiers as variadic arguments
func (b *Builder) WithRecipientIdentifiers(identifiers ...identifier.Identifier) *Builder {
	b.opts.recipients = identifier.NewCollection(identifiers...)
	return b
}

func (b *Builder) WithDefaultMessage(message string) *Builder {
	b.opts.fallbackMessage = message
	return b
}

func (b *Builder) WithMessageForTransport(transportKey common.TransportKey, message string) *Builder {
	b.opts.messagePerTransport[transportKey] = message
	return b
}

func (b *Builder) WithSlackOptions(opts ...slack.MsgOption) *Builder {
	b.opts.slackOpts = opts
	return b
}

func (b *Builder) Build() slack2.RichNotification {
	return &b.opts
}

var _ slack2.RichNotification = &builderOpts{}

func (b *builderOpts) ID() event.ID {
	return b.eventID
}

func (b *builderOpts) Type() event.Type {
	return b.eventType
}

func (b *builderOpts) Recipient() identifier.Collection {
	return b.recipients
}

func (b *builderOpts) Render(key common.TransportKey) string {
	if message, ok := b.messagePerTransport[key]; ok {
		return message
	}

	return b.fallbackMessage
}

func (b *builderOpts) AddRecipients(collection identifier.Collection) {
	b.recipients.Merge(collection)
}

func (b *builderOpts) GetSlackOptions() []slack.MsgOption {
	return b.slackOpts
}
