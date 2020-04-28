/* rpcs3-gameupdater - download helpers */

package main

import (
	"math/rand"
	"time"

	"github.com/cavaliercoder/grab"
)

var client *grab.Client

// this sets up the downloader */
func initDownloader() {
	client = grab.NewClient()
}

/* simple downlaoder */

func downloadFile(folderPath string, url string) (string, error) {
	// Get the data
	req, err := grab.NewRequest(folderPath, url)
	if err != nil {
		return "", err
	}

	resp := client.Do(req)
	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	start := time.Now()

Loop:
	for {
		select {
		case <-t.C:
			sameLinePrint("Transferred %v / %v bytes (%.2f%%) at %.0f Mb/s",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress(),
				resp.BytesPerSecond()/1024/1024)

		case <-resp.Done:
			sameLinePrint("Transferred %v / %v bytes (%.2f%%) at %.0f Mb/s",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress(),
				float64(resp.Size())/(time.Now().Sub(start).Seconds())/1024/1024)
			stopSameLinePrint()
			break Loop
		}
	}
	if err := resp.Err(); err != nil {
		printError("Download failed: %v\n", err)
		return "", err
	}
	return resp.Filename, err
}

/* has logic to retry and sleep */

func downloadFileWithRetries(folderPath string, url string, sha string) {
	for i := 0; i < fetchConfig().DLRetries; i++ {
		time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
		fileName, err := downloadFile(folderPath, url)
		if err != nil {
			printError("Couldn't download '%s' at '%s' (errorcode: '%s')", url, fileName, err)
			continue
		}
		if verifyChecksums(fileName, sha) {
			return
		}
	}
}
