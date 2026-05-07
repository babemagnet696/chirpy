package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

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

func TestCreateJWT(t *testing.T) {
	// Arrange
	tokenSecret :="test-random-long-variable"
	expiresIn := 72 * time.Hour
	userID := uuid.New()
	// Act
	signedString, err := MakeJWT(userID, tokenSecret, expiresIn)
	

	// Assert
	if err != nil {
		t.Fatalf("Error creating signed string: %s", err.Error())
	}
	if signedString == "" {
		t.Errorf("Signed string is empty")
	}
}

func TestValidateJWT(t *testing.T) {
	// Arrange
	tokenSecret := "test-random-long-variable"
	expiresIn   := 72 * time.Hour
	userID      := uuid.New()
	signedString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Error creating signed string: %s", err.Error())
	}

	// Act
	id_correct, err_c := ValidateJWT(signedString, tokenSecret)
	id_incorrect, err_i := ValidateJWT(signedString, "fake-secret")
	
	// Assert
	if err_c != nil {
		t.Fatalf("Error validating JWT: %s", err_c.Error())
	}
	if err_i == nil {
		t.Fatalf("JWT validated when it should have failed")
	}
	if id_correct != userID {
		t.Errorf("validated id does not match userID: %v != %v", id_correct, userID)
	}
	if id_incorrect == userID {
		t.Errorf("Invalidated id matches userID")
	}

}