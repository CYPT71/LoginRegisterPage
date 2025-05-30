package utils

import "testing"

func TestPartialUserUnmarshal(t *testing.T) {
	jsonData := []byte(`{"icon":"i","email":"e","password":"p"}`)
	var u PartialUser
	if err := u.Unmarshal(jsonData); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Icon != "i" || u.Email != "e" || u.Password != "p" {
		t.Fatalf("unmarshal did not populate fields")
	}
}
