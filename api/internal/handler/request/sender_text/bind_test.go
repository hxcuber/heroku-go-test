package sender_text

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRequest_Bind(t *testing.T) {
	type test struct {
		request *Request
		expErr  error
	}

	for s, tc := range map[string]test{
		"empty": {
			request: &Request{
				Sender: "",
			},
			expErr: errors.New("sender is a required field"),
		},
		"invalid_email": {
			request: &Request{
				Sender: "test",
			},
			expErr: errors.New("sender field must be an email"),
		},
		"success": {
			request: &Request{
				Sender: "test@test.com",
			},
			expErr: nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			if tc.expErr != nil {
				require.ErrorContains(t, tc.request.Bind(nil), tc.expErr.Error())
			} else {
				require.ErrorIs(t, tc.request.Bind(nil), tc.expErr)
			}
		})
	}
}
