package main

import (
	// TODO: these are the UI libs
	//"fyne.io/fyne/app"
	//"fyne.io/fyne/widget"
	// TODO: figure out if we really need this

	"bufio"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/* constants */

const urlPattern string = "https://a0.ww.np.dl.playstation.net/tpl/np/%s/%s-ver.xml"
const gamesYAML = "games.yml"

/* structures */

/* this is the sony titlepatch format */

// Paramsfo contains the title of the game
type Paramsfo struct {
	TITLE string `xml:"TITLE"`
}

// Package contains the actual patch file along with metadata (size, sha1)
type Package struct {
	Size     string   `xml:"size,attr"`
	SHA1     string   `xml:"sha1sum,attr"`
	Paramsfo Paramsfo `xml:"paramsfo"`
	URL      string   `xml:"url,attr"`
	Version  string   `xml:"version,attr"`
}

// Tag wraps packages
type Tag struct {
	// we might have multiple patches in a title
	Package []Package `xml:"package"`
	Name    string    `xml:"name,attr"`
	Popup   string    `xml:"popup,attr"`
	Signoff string    `xml:"signoff,attr"`
}

// TitlePatch contains all of the patch data for a given title
type TitlePatch struct {
	Tag     Tag    `xml:"tag"`
	Status  string `xml:"status,attr"`
	Titleid string `xml:"titleid,attr"`
}

// just a generic helper

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

/* parses the given games.yml file and returns a list of update urls */

func parseGamesToURLs(gamesFile string) []string {
	fmt.Printf("Parsing '%s'\n", gamesFile)

	var urlList []string
	file, err := os.Open(gamesFile)

	if isError(err) {
		fmt.Printf("Couldn't open '%s' (errorcode: %d)", gamesFile, err)
		return urlList
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		gameID := strings.Split(line, ":")[0]
		url := fmt.Sprintf(urlPattern, gameID, gameID)

		fmt.Printf("found game: '%s', url:'%s'\n", gameID, url)
		urlList = append(urlList, url)

		// we use err to figure out end of input
		if isError(err) {
			return urlList
		}
	}
}

func main() {
	fmt.Println("downloading using games.yml")

	urls := parseGamesToURLs(gamesYAML)

	for index, url := range urls {
		fmt.Printf("fetching URL %d: '%s'\n", index, url)

		// we need this because we can't verify the signature
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpClient := &http.Client{Transport: transport}
		response, err := httpClient.Get(url)

		if isError(err) {
			fmt.Printf("Error: Can't open url '%s'\n", url)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)

		if isError(err) {
			fmt.Printf("can't read response body\n")
			break
		}
		patch := TitlePatch{}
		err = xml.Unmarshal([]byte(body), &patch)

		if isError(err) {
			fmt.Printf("can't parse response XML\n")
			break
		}

		fmt.Printf("title '%s (%s) url '%s'\n",
			patch.Tag.Package[0].Paramsfo.TITLE,
			patch.Titleid,
			patch.Tag.Package[0].URL)
	}

	// TODO: UI stuff
	//app := app.New()

	//w := app.NewWindow("Hello")
	//w.SetContent(widget.NewLabel("Hello Fyne!"))

	//w.ShowAndRun()

}
