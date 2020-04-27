/* rpcs3-gameupdater - download helpers */

package main

import (
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
	for i := 0; i < fetchConfig().DLRetries; i++ {
		time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
		err := downloadFile(filePath, url)
		if err != nil {
			printError("Couldn't download '%s' at '%s' (errorcode: '%s')", url, filePath, err)
			continue
		}
		if verifyChecksums(filePath, sha) {
			return
		}
	}
}
