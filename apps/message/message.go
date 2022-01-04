package message

import (
	"github.com/getevo/evo-ng"
	"time"
)

const (
	ERROR   Type = "error"
	SUCCESS Type = "success"
	WARNING Type = "warning"
	INFO    Type = "info"
)

// Type helper
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

// Message creates flash message
//  @receiver c
//  @param typ
//  @param message
func (c *Context) Message(typ Type, message string) {
	if len(c.messages) == 0 {
		var err = c.request.Cookie("message").ParseJSON(&c.messages)
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

// GetMessages returns list of flash messages
//  @receiver c
//  @return []Message
func (c *Context) GetMessages() []Message {
	if len(c.messages) == 0 {
		c.request.Cookie("message").ParseJSON(&c.messages)
	}
	return c.messages
}

// Flush flash messages
//  @receiver c
//  @return []Message
func (c *Context) Flush() []Message {
	c.request.ClearCookie("message")
	c.messages = []Message{}
	return c.messages
}

// Extend parse request
//  @receiver context
//  @param request
//  @return error
func (context *Context) Extend(request *evo.Context) error {
	context.request = request
	return nil
}

// Error set flash error message
//  @receiver c
//  @param message
func (c *Context) Error(message string) {
	c.Message(ERROR, message)
}

// Info set flash info message
//  @receiver c
//  @param message
func (c *Context) Info(message string) {
	c.Message(INFO, message)
}

// Warning set flash warning message
//  @receiver c
//  @param message
func (c *Context) Warning(message string) {
	c.Message(WARNING, message)
}

// Success set flash success message
//  @receiver c
//  @param message
func (c *Context) Success(message string) {
	c.Message(SUCCESS, message)
}
