/* rpcs3-gameupdater - utilities */

package main

import (
	"encoding/hex"
	"os"
)

// check for file or path existence
func testPath(path string) bool {
	return true
}

// IsTTY checks if we're in a terminal, this is platform specific, implementations in _unix and _windows files
func IsTTY() bool {
	return isTTY()
}

/* verify that the 3 checksums (passed, stored and calculated) match */
func verifyPKGChecksums(filePath string, sha string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		printError("Couldn't open '%s' (errorcode: '%s')\n", filePath, err)
		return false
	}
	stat, err := file.Stat()
	if err != nil {
		printError("Couldn't Stat for file '%s' (errorcode: '%s')\n", filePath, err)
		return false
	}
	// skip reading last 0x20 per
	// https://github.com/13xforever/psn-pkg-validator
	// those bytes contain the csum
	sizeToRead := stat.Size() - 0x20
	computedSHA, err := computeSHA1N(file, sizeToRead)
	if computedSHA == "" || err != nil {
		printError("Couldn't compute SHA1 for file '%s' (errorcode: '%s')\n", filePath, err)
		return false
	}
	buf := make([]byte, 0x20)
	file.Read(buf)
	storedSHA := hex.EncodeToString(buf[:20])
	if storedSHA == "" {
		printError("Couldn't retrieve stored SHA1 for file '%s' (errorcode: '%s')\n", filePath)
		return false
	}
	printDebug("passedSHA: %s", sha)
	printDebug("calcedSHA: %s", computedSHA)
	printDebug("storedSHA: %s", storedSHA)
	if sha == computedSHA {
		printDebug("The passed sha1 matches the computed one for '%s'", filePath)
		if storedSHA == computedSHA {
			printDebug("The stored sha1 matches the computed one for '%s'", filePath)
			return true
		}
		printDebug("The stored sha1 does not match the computed one for '%s'", filePath)
		return false

	}
	printDebug("The passed sha1 does not match the computed one for '%s'", filePath)
	return false
}
