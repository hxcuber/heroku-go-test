package requestor_target

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
		"empty_requestor": {
			request: &Request{
				Requestor: "",
				Target:    "",
			},
			expErr: errors.New("requestor is a required field"),
		},
		"empty_target": {
			request: &Request{
				Requestor: "req@test",
				Target:    "",
			},
			expErr: errors.New("target is a required field"),
		},
		"invalid_req": {
			request: &Request{
				Requestor: "req@test",
				Target:    "tar@test",
			},
			expErr: errors.New("requestor field must be an email"),
		},
		"invalid_tar": {
			request: &Request{
				Requestor: "req@test.com",
				Target:    "tar@test",
			},
			expErr: errors.New("target field must be an email"),
		},
		"duplicate": {
			request: &Request{
				Requestor: "test@test.com",
				Target:    "test@test.com",
			},
			expErr: errors.New("emails cannot be the same"),
		},
		"success": {
			request: &Request{
				Requestor: "req@test.com",
				Target:    "tar@test.com",
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
