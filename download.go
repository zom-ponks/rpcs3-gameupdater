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
		if _, err := io.Copy(hash, file); err != nil {
			printError("Couldn't copy the file data for '%s' to the sha (errorcode: '%s')", filePath, err)
			continue
		}
		hashInBytes := hash.Sum(nil)[:20]
		SHA1String := hex.EncodeToString(hashInBytes)
		printDebug("passed sha: '%s'", sha)
		printDebug("calculated sha: '%s'", SHA1String)
		if sha == SHA1String {
			printDebug("The sha1 matches for url '%s' at '%s'", url, filePath)
			break
		} else {
			printDebug("The sha1 does not match for url '%s' at '%s'", url, filePath)
		}
	}

}
