package creds_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"

	"github.com/derpartizanen/gophkeeper/internal/libraries/creds"
)

func TestPasswordToString(t *testing.T) {
	tt := []struct {
		name string
		data string
	}{
		{
			name: "Print password",
			data: "1q2w3e",
		},
		{
			name: "Print empty password",
			data: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sat := creds.Password(tc.data)
			snaps.MatchSnapshot(t, sat.String())
		})
	}
}

func TestConnURIToString(t *testing.T) {
	tt := []struct {
		name     string
		data     string
		expected string
	}{
		{
			name: "Print database URI",
			data: "postgres://postgres:postgres@127.0.0.1:5432/goph?sslmode=disable",
		},
		{
			name: "Print empty secret",
			data: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sat := creds.ConnURI(tc.data)

			snaps.MatchSnapshot(t, sat.String())
		})
	}
}
