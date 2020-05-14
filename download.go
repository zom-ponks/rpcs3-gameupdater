/* rpcs3-gameupdater - download helpers */

package main

import (
	"crypto/tls"
	"io/ioutil"
	"path/filepath"
	//"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
)

var client *grab.Client

// this sets up the downloader */
func initDownloader() {
	client = grab.NewClient()
	// we need this because we can't verify the signature
	client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.HTTPClient.Timeout = time.Duration(conf.DLTimeout) * time.Second
}

/* simple downloader that supports a UI and resuming */

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
			sameLinePrint("Transferred %v/%v Mb (%.2f%%) at %.0f Mb/s",
				resp.BytesComplete()/1024/1024,
				resp.Size/1024/1024,
				100*resp.Progress(),
				resp.BytesPerSecond()/1024/1024)

		case <-resp.Done:
			sameLinePrint("Transferred %v/%v Mb (%.2f%%) at %.0f Mb/s",
				resp.BytesComplete()/1024/1024,
				resp.Size/1024/1024,
				100*resp.Progress(),
				float64(resp.Size)/(time.Now().Sub(start).Seconds())/1024/1024)
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

/* removes ".tmp" from the filename */

func renameToNoTmp(fileName string) string {
	newName := strings.Replace(fileName, ".tmp", "", 1)
	os.Rename(fileName, newName)
	return newName
}

// VerifyChecksums is a function passed to downloadFileWithRetries
// to verify the PKG after downloading it
type VerifyChecksums func(string, string) bool

/* has logic to retry and sleep */

func downloadWithRetries(folderPath string, url string, sha string, verifyChecksums VerifyChecksums) string {
	filePath := path.Join(folderPath, filepath.Base(url))
	if pathExists(filePath) {
		return filePath
	}
	filePath = filePath + ".tmp"
	for i := 0; i < fetchConfig().DLRetries; i++ {
		//time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
		fileName, err := downloadFile(filePath, url)
		if err != nil {
			printError("Couldn't download '%s' at '%s' (errorcode: '%s')", url, fileName, err)
			continue
		}
		if verifyChecksums != nil {
			if verifyChecksums(fileName, sha) {
				return renameToNoTmp(fileName)
			}
			// because the downloader retries based on file size
			// keeping the wrong file around could be a problem
			// it may be better to keep it for a retry or so first though
			err := os.Remove(fileName)
			if err != nil {
				printError("Couldn't delete file '%s' (errorcode: '%s')", fileName, err)
			}
			continue
		}
		return renameToNoTmp(fileName)
	}
	return ""
}

/* download the XML */

func getXML(url string) []byte {
	fileName := downloadWithRetries(fetchConfig().XMLCachePath, url, "", nil)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		printError("Error reading the file at %s (errorcode: %s)", fileName, err)
	}
	return data

}

/* download the PKG file */

func getPKG(url string, sha string) bool {
	fileName := downloadWithRetries(fetchConfig().PkgDLPath, url, sha, verifyPKGChecksums)
	if fileName != "" {
		return false
	}
	return true
}
