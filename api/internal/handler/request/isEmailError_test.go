package request

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsEmailError(t *testing.T) {
	type test struct {
		email  string
		expErr error
	}

	for s, tc := range map[string]test{
		"error": {
			email:  "test",
			expErr: errors.New("email field must be an email"),
		},
		"success": {
			email:  "test@test.com",
			expErr: nil,
		},
	} {
		t.Run(s, func(t *testing.T) {
			if tc.expErr != nil {
				require.ErrorContains(t, IsEmailError(tc.email, "email"), tc.expErr.Error())
			} else {
				require.ErrorIs(t, IsEmailError(tc.email, "email"), tc.expErr)
			}
		})
	}
}
