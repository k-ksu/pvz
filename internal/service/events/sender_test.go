package events

import (
	"encoding/json"
	"testing"
	"time"

	"HomeWork_1/internal/model"
	"HomeWork_1/internal/service/events/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestSender_Handle(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	producer := mocks.NewProducer(mc)
	sender := NewSender(false, producer)
	producer.SendSyncMessageMock.Return(0, 0, nil)

	event := model.EventMessage{
		Method:    "listOrders",
		Args:      []string{"--clientID=78", "--action=1"},
		TimeStamp: time.Now(),
	}

	msg, err := json.MarshalIndent(event, "  ", "  ")
	require.NoError(t, err, msg)

	tests := []struct {
		desc      string
		msg       model.EventMessage
		wantError bool
	}{
		{
			desc:      "Test case 1: correct msg",
			msg:       event,
			wantError: false,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			result := sender.SendMessage(&tc.msg)

			if tc.wantError {
				require.Error(t, result)
			} else {
				require.NoError(t, result)
			}
		})
	}
}
