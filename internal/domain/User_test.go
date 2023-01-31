package domain

import "testing"

func TestUserModel_Create(t *testing.T) {
	tests := []struct {
		name string
		user *UserModel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.user.Create()
		})
	}
}
