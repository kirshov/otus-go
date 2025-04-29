package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
	var roleAdmin, roleGuest UserRole
	roleAdmin = "admin"
	roleGuest = "guest"

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123213",
				Name:   "Guest",
				Age:    15,
				Email:  "test@",
				Role:   roleGuest,
				Phones: []string{"12345"},
			},
			expectedErr: errors.New(ValidationErrors{
				ValidationError{Field: "ID", Err: fmt.Errorf(validateStrLenErr, 36)},
				ValidationError{Field: "Age", Err: fmt.Errorf(validateIntMinErr, 18)},
				ValidationError{Field: "Email", Err: errors.New(validateStrRegxErr)},
				ValidationError{Field: "Role", Err: fmt.Errorf(validateInListErr, "admin,stuff")},
				ValidationError{Field: "Phones", Err: fmt.Errorf(validateStrLenErr, 11)},
			}.Error()),
		},
		{
			in: User{
				ID:     "1234567890-1234567890-1234567890-123",
				Name:   "Guest",
				Age:    30,
				Email:  "test@test.test",
				Role:   roleAdmin,
				Phones: []string{"01234567890"},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0",
			},
			expectedErr: errors.New(ValidationErrors{
				ValidationError{Field: "Version", Err: fmt.Errorf(validateStrLenErr, 5)},
			}.Error()),
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    make([]byte, 0),
				Payload:   make([]byte, 0),
				Signature: make([]byte, 0),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 301,
			},
			expectedErr: errors.New(ValidationErrors{
				ValidationError{Field: "Code", Err: fmt.Errorf(validateInListErr, "200,404,500")},
			}.Error()),
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
