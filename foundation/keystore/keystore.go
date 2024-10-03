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
	// keyFileHandle, err := os.Open(KeyPath)
	// if err != nil {
	// 	return "", fmt.Errorf("encryption master key not found: %w", err)
	// }
	//

	// limit Key file size to 1 megabyte. This should be reasonable for
	// almost any key file.
	// keyBytes, err := io.ReadAll(io.LimitReader(keyFileHandle, 1024*1024))
	// if err != nil {
	// 	return "", fmt.Errorf("reading encryption master key: %w", err)
	// }

	return string(keyBytes), nil
}
