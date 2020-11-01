package twist

import (
	"context"
	http "net/http"
	"testing"

	"github.com/dghubble/sling"
	"github.com/go-joe/joe"
	"github.com/go-joe/joe/joetest"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockClient *mockClientAPI

func newTestBotAdapter(t *testing.T) (*BotAdapter, error) {
	ctx := context.Background()
	mockClient = new(mockClientAPI)
	events := make(chan requestCatcher)
	a, err := NewBotAdapter(ctx, events, mockClient)
	require.NoError(t, err)
	return a, err
}

func TestAdapter_Messages(t *testing.T) {
	brain := joetest.NewBrain(t)
	a, _ := newTestBotAdapter(t)

	done := make(chan bool)
	go func() {
		a.handleMessageEvent(brain.Brain)
		done <- true
	}()
	evt := requestCatcher{
		Content:     "Hello World!",
		URLCallback: "twist.test.dev/callback",
		MessageID:   78912,
		UserName:    "test_bot",
	}
	a.events <- evt

	close(a.events)
	<-done
	brain.Finish()

	events := brain.RecordedEvents()
	require.NotEmpty(t, events)
	expectedVlt := joe.ReceiveMessageEvent{
		Text:     "Hello World!",
		Channel:  "twist.test.dev/callback",
		ID:       "78912",
		AuthorID: "test_bot",
	}
	assert.Equal(t, expectedVlt, events[0])
}

func TestAdapter_SendCallback(t *testing.T) {
	a, _ := newTestBotAdapter(t)
	var data requestBodyCallback
	data.Content = "Hello World!"
	req, _ := sling.New().Post("twist.test.dev/callback").BodyJSON(&data).Request()
	expectedRes := &http.Response{}
	mockClient.On("Send", req).Return(expectedRes, nil)
	res, err := a.clientAPI.Send(req)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

type mockClientAPI struct {
	mock.Mock
}

func (_m *mockClientAPI) Send(request *http.Request) (*http.Response, error) {
	ret := _m.Called(request)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(*http.Request) *http.Response); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
