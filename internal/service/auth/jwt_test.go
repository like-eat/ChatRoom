package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseToken(t *testing.T) {
	tokenString, err := GenerateToken("U2026072118994456861", 1)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ParseToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Subject != "U2026072118994456861" || claims.IsAdmin != 1 {
		t.Fatalf("unexpected claims: subject=%q admin=%d", claims.Subject, claims.IsAdmin)
	}
}

func TestRejectsUnexpectedSigningMethod(t *testing.T) {
	claims := Claims{
		TokenUse: accessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "U2026072118994456861",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseToken(tokenString); err == nil {
		t.Fatal("token using alg=none must be rejected")
	}
}
