/* rpcs3-gameupdater - crypto functions */

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

/* computes the SHA1 checksum from a file descriptor up to n bytes */

func computeSHA1N(file *os.File, n int64) (string, error) {
	// Would be faster to do that somewhere else
	hash := sha1.New()

	if _, err := io.CopyN(hash, file, n); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:20]
	return hex.EncodeToString(hashInBytes), nil
}
