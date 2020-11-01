package twist

import (
	"context"
	"fmt"
	"os"

	"github.com/dghubble/sling"
	"github.com/go-joe/joe"
	"github.com/gofiber/fiber/v2"
)

var defaultHttpPort = ":80"

// Catch all data from twist when hit our endpoint
type requestCatcher struct {
	WorkspaceID       int    `json:"workspace_id" form:"workspace_id" xml:"workspace_id"`
	URLCallback       string `json:"url_callback" form:"url_callback" xml:"url_callback"`
	UserID            int    `json:"user_id" form:"user_id" xml:"user_id"`
	IntegrationType   string `json:"integration_type" form:"integration_type" xml:"integration_type"`
	ConversationTitle string `json:"conversation_title" form:"conversation_title" xml:"conversation_title"`
	MessageID         int    `json:"message_id" form:"message_id" xml:"message_id"`
	Content           string `json:"content" form:"content" xml:"content"`
	URLTTL            int    `json:"url_ttl" form:"url_ttl" xml:"url_ttl"`
	UserName          string `json:"user_name" form:"user_name" xml:"user_name"`
	VerifyToken       string `json:"verify_token" form:"verify_token" xml:"verify_token"`
	EventType         string `json:"event_type" form:"event_type" xml:"event_type"`
}

type requestBodyCallback struct {
	Content string `json:"content"`
}

type BotAdapter struct {
	context   context.Context
	events    chan requestCatcher
	clientAPI clientAPI
	app       *fiber.App
}

func NewBotAdapter(ctx context.Context, events chan requestCatcher, clientApi clientAPI) (*BotAdapter, error) {
	return &BotAdapter{
		context:   ctx,
		events:    events,
		clientAPI: clientApi,
		app:       fiber.New(),
	}, nil
}

func newAdapter(ctx context.Context) (*BotAdapter, error) {
	events := make(chan requestCatcher)
	a, _ := NewBotAdapter(ctx, events, &httpClient{})
	go func() {
		a.app.Post("/messages", func(c *fiber.Ctx) error {
			var p requestCatcher
			if err := c.BodyParser(&p); err != nil {
				return err
			}
			events <- p
			c.Status(204)
			return nil
		})
		portListen := defaultHttpPort
		// just set on environment variable if want to change where service will running on port number
		// default to 80
		if os.Getenv("HTTP_LISTEN_PORT") != "" {
			portListen = ":" + os.Getenv("HTTP_LISTEN_PORT")
		}
		// since twist bot integration using outgoing webhook,
		// this service will listen to port as http
		a.app.Listen(portListen)
	}()
	return a, nil
}

// Adapter returns a new BotAdapter as joe.Module.
func Adapter() joe.Module {
	return joe.ModuleFunc(func(joeConf *joe.Config) error {
		a, err := newAdapter(joeConf.Context)
		if err != nil {
			return err
		}
		joeConf.SetAdapter(a)
		return nil
	})
}

func (b *BotAdapter) handleMessageEvent(brain *joe.Brain) {
	for req := range b.events {
		brain.Emit(joe.ReceiveMessageEvent{
			Text:     req.Content,
			Channel:  req.URLCallback,           // parse url callback from twist to channel, need to send response to message
			ID:       fmt.Sprint(req.MessageID), // since ID on joe required string, need to convert from integer to string
			AuthorID: req.UserName,
		})
	}
}

func (b *BotAdapter) RegisterAt(brain *joe.Brain) {
	go b.handleMessageEvent(brain)
}

func (b *BotAdapter) Send(text, urlCallback string) error {
	var data requestBodyCallback
	data.Content = text
	twistReq, _ := sling.New().Post(urlCallback).BodyJSON(&data).Request()
	_, err := b.clientAPI.Send(twistReq)
	return err
}

func (b *BotAdapter) Close() error {
	defer close(b.events)
	return nil
}
