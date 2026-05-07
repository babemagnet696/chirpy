package auth

import "testing"

func TestPasswordHash(t *testing.T) {
	// Arrange
	password := "BlackAndBlue123"

	// Act
	hash, err := HashPassword(password)

	// Assert
	if err != nil {
		t.Fatalf("Error creating hash: %s", err.Error())
	}
	if hash == "" {
		t.Errorf("Hash is empty")
	}
	if hash == password {
		t.Errorf("Hash is the same as password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// Arrange
	password := "BlackAndBlue123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error creating hash: %s", err.Error())
	}

	// Act
	ok, err := CheckPasswordHash(password, hash)

	// Assert
	if err != nil {
		t.Fatalf("Error checking hash: %s", err.Error())
	}

	if !ok {
		t.Errorf("Passwords should match but dont")
	}
}