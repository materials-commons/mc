package pw

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"testing"

	"golang.org/x/crypto/pbkdf2"
)

func TestPW(t *testing.T) {
	originalSalt := "682ea42af2c411e6a622005056ab2d9f"
	salt := fmt.Sprintf("$p5k2$%x$%s", 4000, originalSalt)
	fmt.Printf("Using salt: '%s'\n", salt)
	expected := "B6j2HHO7fNboyqHutqxVIcKIQUIIrt9L"
	what := pbkdf2.Key([]byte("Dr0gan5ST!"), []byte(salt), 4000, 24, sha1.New)
	got := base64.StdEncoding.EncodeToString(what)
	fmt.Printf("%s:%d:%s:%d\n", got, len(got), expected, len(expected))
	//B6j2HHO7fNboyqHutqxVIcKIQUIIrt9L
}
