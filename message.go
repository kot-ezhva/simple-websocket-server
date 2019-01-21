package main

import (
	"encoding/json"
)

type Message struct {
	Event             string `json:"event,omitempty"` // subscribe, newMessage, typing
	UserId            int    `json:"userId,omitempty"`
	RelatedObjectType string `json:"relatedObjectType,omitempty"` // customer_case, shophelp, social
	RelatedObjectId   int    `json:"relatedObjectId,omitempty"`
	NewMessageId      int    `json:"newMessageId,omitempty"`
	User              struct {
		Avatar   string `json:"avatar,omitempty"`
		UserName string `json:"userName,omitempty"`
	} `json:"user,omitempty"`
}

func parseMessage(msg []byte) *Message {
	message := Message{}
	json.Unmarshal([]byte(msg), &message)
	return &message
}

func (message *Message) isSubscribe() bool {
	return message.Event == "subscribe"
}

func (message *Message) toSend() []byte {
	msg, _ := json.Marshal(message)
	return msg
}
