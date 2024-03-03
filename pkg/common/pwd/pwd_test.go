package pwd

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "a12345"
	hashedPassword := HashPassword(password)
	t.Log("hashed: ", hashedPassword)
	if !CheckPasswordHash(hashedPassword, password) {
		t.Errorf("CheckPasswordHash failed for password: %s and hash: %s", password, hashedPassword)
	}

	wrongPassword := "w2w2w2"
	if CheckPasswordHash(hashedPassword, wrongPassword) {
		t.Errorf("CheckPasswordHash should fail for password: %s and hash: %s", wrongPassword, hashedPassword)
	}
}
