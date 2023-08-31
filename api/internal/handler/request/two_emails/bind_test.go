package two_emails

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRequest_Bind(t *testing.T) {
	type test struct {
		request *Request
		expErr  error
	}

	for s, tc := range map[string]test{
		"fail_no_emails": {
			request: &Request{
				Friends: nil,
			},
			expErr: errors.New("friends is a required field"),
		},
		"fail_1_email": {
			request: &Request{
				Friends: []string{"test1@test.com"},
			},
			expErr: errors.New("2 elements required, less than 2 given"),
		},
		"fail_4_emails": {
			request: &Request{
				Friends: []string{"test1@test.com", "test2@test.com", "test3@test.com", "test4@test.com"},
			},
			expErr: errors.New("2 elements required, more than 2 given"),
		},
		"fail_invalid_emails": {
			request: &Request{
				Friends: []string{"test", "hello"},
			},
			expErr: errors.New("elements of friends field must be an email"),
		},
		"fail_duplicate_emails": {
			request: &Request{
				Friends: []string{"test1@test.com", "test1@test.com"},
			},
			expErr: errors.New("emails cannot be the same"),
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
