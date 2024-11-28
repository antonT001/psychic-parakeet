package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/antonT001/psychic-parakeet/hw09_struct_validator/e"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid",
			in: User{
				ID:    "012345678901234567890123456789012345",
				Name:  "Anton",
				Age:   40,
				Email: "anton@example.com",
				Role:  "admin",
				Phones: []string{
					"89999999999",
					"89111111111",
				},
				meta: []byte("qwertyui"),
			},
			expectedErr: nil,
		},
		{
			name: "invalid len",
			in: App{
				Version: "123456",
			},
			expectedErr: e.VError("Version", errors.New("the value does not match the specified length")),
		},
		{
			name: "no tag and unsupported value",
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			name: "not included integer and not validate tag",
			in: Response{
				Code: 1,
				Body: "body",
			},
			expectedErr: e.VError("Code", errors.New("the value is not included in the specified set")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			t.Log(err)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
