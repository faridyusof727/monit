package utils

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func FirebaseAuth(credFile string) *auth.Client {
	opt := option.WithCredentialsFile(credFile)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	return auth
}
