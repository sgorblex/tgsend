package tdwrap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"golang.org/x/crypto/ssh/terminal"
)

type termAuth struct{}

func GetFlow() auth.Flow {
	return auth.NewFlow(
		termAuth{},
		auth.SendCodeOptions{},
	)
}

func (c termAuth) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c termAuth) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

func (a termAuth) Phone(_ context.Context) (string, error) {
	fmt.Print("Enter phone number with prefix: ")
	var phone string
	_, err := fmt.Scan(&phone)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(phone), nil
}

func (a termAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (a termAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	var code string
	_, err := fmt.Scan(&code)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}
