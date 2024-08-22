package events

import (
	"encoding/json"
	"testing"
	"time"

	"HomeWork_1/internal/model"

	"github.com/stretchr/testify/require"
)

func TestReader_Handle(t *testing.T) {
	t.Parallel()

	reader := NewReader(false)
	event := model.EventMessage{
		Method:    "listOrders",
		Args:      []string{"--clientID=78", "--action=1"},
		TimeStamp: time.Now(),
	}

	msg, err := json.MarshalIndent(event, "  ", "  ")
	require.NoError(t, err, msg)

	tests := []struct {
		desc      string
		msg       []byte
		wantError bool
	}{
		{
			desc:      "Test case 1: correct msg",
			msg:       msg,
			wantError: false,
		},
		{
			desc:      "Test case 2: not correct msg",
			msg:       []byte{0},
			wantError: true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			result := reader.Handle(tc.msg)

			if tc.wantError {
				require.Error(t, result)
			} else {
				require.NoError(t, result)
			}
		})
	}
}
