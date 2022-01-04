package message

import (
	"encoding/json"
	"github.com/getevo/evo-ng"
	"time"
)

const (
	ERROR   Type = "error"
	SUCCESS Type = "success"
	WARNING Type = "warning"
	INFO    Type = "info"
)

type Type string

// Message contains message type and text
type Message struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

// Context extend Message context
type Context struct {
	request  *evo.Context
	messages []Message
}

func (c *Context) Message(typ Type, message string) {
	if len(c.messages) == 0 {
		var err = json.Unmarshal([]byte(c.request.Cookie("message")), &c.messages)
		if err != nil {
			c.messages = []Message{}
		}
	}
	var found = false
	for _, item := range c.messages {
		if typ == item.Type && item.Message == message {
			found = true
			break
		}
	}
	if !found {
		c.messages = append(c.messages, Message{typ, message})
		c.request.Cookie("message", c.messages, time.Duration(60*time.Second))
	}
}

func (c *Context) GetMessages() []Message {
	if len(c.messages) == 0 {
		json.Unmarshal([]byte(c.request.Cookie("message")), &c.messages)
	}
	return c.messages
}

func (c *Context) Flush() []Message {
	c.request.ClearCookie("message")
	c.messages = []Message{}
	return c.messages
}

func (context *Context) Extend(request *evo.Context) error {
	context.request = request
	return nil
}

func (c *Context) Error(message string) {
	c.Message(ERROR, message)
}

func (c *Context) Info(message string) {
	c.Message(INFO, message)
}

func (c *Context) Warning(message string) {
	c.Message(WARNING, message)
}

func (c *Context) Success(message string) {
	c.Message(SUCCESS, message)
}
