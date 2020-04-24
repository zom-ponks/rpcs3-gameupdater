/* rpcs3-gameupdater - downlaod helpers */

package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

/* simple downlaoder */

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	file, err := os.Create(filepath)
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

func downloadFileWithRetries(filepath string, url string) {
	time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
	for i := 0; i < fetchConfig().DLRetries; i++ {
		err := downloadFile(filepath, url)
		if err == nil {
			break
		}
	}

}
