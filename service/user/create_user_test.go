package user_service

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestCreateUserValidation(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	tests := []struct {
		name    string
		request CreateUserRequest
		wantErr bool
	}{
		{"ข้อมูลถูกต้อง", CreateUserRequest{Username: "thanutsu", Name: "thanut"}, false},
		{"username สั้นกว่า 3", CreateUserRequest{Username: "ab", Name: "thanut"}, true},
		{"username ว่าง", CreateUserRequest{Username: "", Name: "thanut"}, true},
		{"name ว่าง", CreateUserRequest{Username: "thanutsu", Name: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
