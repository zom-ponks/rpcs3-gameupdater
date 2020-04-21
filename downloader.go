package main

import (
	//"fyne.io/fyne/app"
	//"fyne.io/fyne/widget"
	//"gopkg.in/yaml.v2"

	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	//"log"
)

/* constants */

const urlPattern string = "https://a0.ww.np.dl.playstation.net/tpl/np/%s/%s-ver.xml"
const gamesYAML = "games.yml"

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
		}
		fmt.Printf("Response body: \n'''\n%s\n'''\n", body)

	}

	// TODO: UI stuff
	//app := app.New()

	//w := app.NewWindow("Hello")
	//w.SetContent(widget.NewLabel("Hello Fyne!"))

	//w.ShowAndRun()

}
