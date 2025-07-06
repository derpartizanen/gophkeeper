// Package creds provides common types for wrapping sensitive data in the application.
package creds

import (
	"regexp"
	"strings"
)

// Password is sensitive value which shouldn't leak to logs.
type Password string

// String converts password to string.
func (p Password) String() string {
	return strings.Repeat("*", len(p))
}

// ConnURI is URI with sensitive values (e.g. login:password) which shouldn't leak to logs.
type ConnURI string

var _URISecrets = regexp.MustCompile(`(://).*:.*(@)`)

// String converts SecretURI to string.
func (u ConnURI) String() string {
	return string(_URISecrets.ReplaceAll([]byte(u), []byte("$1*****:*****$2")))
}
