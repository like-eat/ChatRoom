package auth

import "testing"

func TestPasswordHashAndVerify(t *testing.T) {
	password := "SafePassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if hash == password {
		t.Fatal("password hash must not equal plaintext")
	}

	valid, needsUpgrade, err := VerifyPassword(hash, password)
	if err != nil {
		t.Fatal(err)
	}
	if !valid || needsUpgrade {
		t.Fatalf("expected bcrypt password to be valid without upgrade, got valid=%v upgrade=%v", valid, needsUpgrade)
	}

	valid, _, err = VerifyPassword(hash, "wrong-password")
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("wrong password must not validate")
	}
}

func TestLegacyPlaintextPasswordRequestsUpgrade(t *testing.T) {
	valid, needsUpgrade, err := VerifyPassword("123456", "123456")
	if err != nil {
		t.Fatal(err)
	}
	if !valid || !needsUpgrade {
		t.Fatalf("expected legacy password migration, got valid=%v upgrade=%v", valid, needsUpgrade)
	}
}
