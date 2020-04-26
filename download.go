/* rpcs3-gameupdater - downlaod helpers */

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

/* simple downlaoder */

func downloadFile(filePath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(file, resp.Body)

	// Supposedly it is bad to defer for writes
	file.Close()
	return err
}

/* has logic to retry and sleep */

func downloadFileWithRetries(filePath string, url string, sha string) {
	hash := sha1.New()
	for i := 0; i < fetchConfig().DLRetries; i++ {
		time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
		err := downloadFile(filePath, url)
		if err != nil {
			printError("Couldn't download '%s' at '%s' (errorcode: '%s')", url, filePath, err)
			continue
		}
		file, err := os.Open(filePath)
		if err != nil {
			printError("Couldn't open '%s' (errorcode: '%s')\n", filePath, err)
			continue
		}
		stat, err := file.Stat()
		if err != nil {
			printError("Couldn't get stat for file '%s' (errorcode: '%s')\n", filePath, err)
			continue
		}
		// skip reading last 0x20 per
		// https://github.com/13xforever/psn-pkg-validator
		// those bytes contain the csum
		sizeToRead := stat.Size() - 0x20
		if _, err := io.CopyN(hash, file, sizeToRead); err != nil {
			printError("Couldn't copy the file data for '%s' to the sha (errorcode: '%s')", filePath, err)
			continue
		}
		hashInBytes := hash.Sum(nil)[:20]
		computedSHA := hex.EncodeToString(hashInBytes)
		buf := make([]byte, 0x20)
		n, _ := file.Read(buf)
		printDebug("bytes copied: " + string(n))
		storedSHA := hex.EncodeToString(buf[:20])
		if sha == computedSHA {
			printDebug("The passed sha1 matches the computed one for url '%s' at '%s'", url, filePath)
			if storedSHA == computedSHA {
				printDebug("The stored sha1 matches the computed one for url '%s' at '%s'", url, filePath)
				break
			} else {
				printDebug("The stored sha1 does not match the computed one for url '%s' at '%s'", url, filePath)
			}

		} else {
			printDebug("The sha1 does not match for url '%s' at '%s'", url, filePath)
		}
	}

}
