package keystore

import (
	_ "embed"
)

var (
	UserDEKeyAlias string
	KeyPath        string
)

//go:embed local-master-key.txt
var keyBytes []byte

func MustGetMKey() (string, error) {

	return string(keyBytes), nil
}
