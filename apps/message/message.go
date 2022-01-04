package message

import (
	"encoding/json"
	"github.com/getevo/evo-ng"
)

const (
	ERROR   Type = "error"
	SUCCESS Type = "success"
	WARNING Type = "info"
	INFO    Type = "info"
)

type Type string

// Message contains message type and text
type Message struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

// Context Extend Message Context
type Context struct {
	request *evo.Context
	flash   []Message
}

func (c *Context) Message(typ Type, message string) {
	if len(c.flash) == 0 {
		json.Unmarshal([]byte(c.request.Cookie("message")), &c.flash)
	}
	c.flash = append(c.flash, Message{typ, message})
	c.request.Cookie("message", c.flash)
}

func (c *Context) GetMessages() []Message {
	if len(c.flash) == 0 {
		json.Unmarshal([]byte(c.request.Cookie("message")), &c.flash)
	}
	return c.flash
}

func (c *Context) Flush() []Message {
	c.request.Cookie("message", "[]")
	c.flash = []Message{}
	return c.flash
}

func (context *Context) Extend(request *evo.Context) error {
	context.request = request
	return nil
}
