package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsEmail(t *testing.T) {
	type test struct {
		email  string
		expRes bool
	}

	for s, tc := range map[string]test{
		"fail_1": {
			email:  "test",
			expRes: false,
		},
		"fail_2": {
			email:  "test@",
			expRes: false,
		},
		"fail_3": {
			email:  "test@test",
			expRes: false,
		},
		"fail_4": {
			email:  "test@test.",
			expRes: false,
		},
		"success_1": {
			email:  "test@test.com",
			expRes: true,
		},
		"success_2": {
			email:  "test.hello@test.com",
			expRes: true,
		},
	} {
		t.Run(s, func(t *testing.T) {
			if tc.expRes {
				require.True(t, IsEmail(tc.email))
			} else {
				require.False(t, IsEmail(tc.email))
			}
		})
	}
}
