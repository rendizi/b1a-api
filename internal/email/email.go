package email

import (
	"errors"

	emailVerifier "github.com/AfterShip/email-verifier"
)

func Verify(email string) error {
	verifier := emailVerifier.NewVerifier()
	ret, err := verifier.Verify(email)
	if err != nil {
		return err
	}
	if !ret.Syntax.Valid {
		return errors.New("syntax is invalid")
	}
	return nil
}
