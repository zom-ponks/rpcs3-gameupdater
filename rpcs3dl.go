/* rpcs3-gameupdater main entry point */

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
	"path/filepath"
	"strings"
	"time"
)

// just a generic helper

func isError(err error) bool {
	if err != nil {
		printError(err.Error())
	}
	return (err != nil)
}

/* parses the given config.yml file and returns the path to dev_hdd0 */

func getGamesPath(configYML string) string {
	printInfo("Parsing '%s'\n", configYML)
	path := "test"
	file, err := os.Open(configYML)

	if isError(err) {
		printError("Couldn't open '%s' (errorcode: %d)", configYML, err)
		return path
	}

	defer file.Close()

	emulatorDir := ""
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if strings.Contains(line, "$(EmulatorDir):") {
			emulatorDir = strings.TrimSpace(strings.Split(line, ":")[1])
			if emulatorDir == "\"\"" {
				emulatorDir = filepath.Dir(configYML) + "/"
			}
			printDebug("emudir: " + emulatorDir + "TT\n")
		}
		if strings.Contains(line, "/dev_hdd0/") {
			path = strings.Replace(strings.TrimSpace(strings.Split(line, ":")[1]), "$(EmulatorDir)", emulatorDir, -1)
			printDebug("path: " + path + "\n")
		}

		// we use err to figure out end of input
		if isError(err) {
			return path
		}
	}
}

/* replaces the variable in the URL with the gameID */

func getURLFromID(id string) string {
	return fmt.Sprintf(urlPattern, id, id)
}

/* gets games URLs from the various folders */

func getGamesURLs(path string) []string {
	var urlList []string

	// first reads the disc folder
	files, err := ioutil.ReadDir(path + "disc")
	if err != nil {
		printError("Couldn't open '%s' (errorcode: %d)", path, err)
		return urlList
	}
	// then reads the game folder
	files2, err := ioutil.ReadDir(path + "game")
	if err != nil {
		printError("Couldn't open '%s' (errorcode: %d)", path, err)
		return urlList
	}
	files = append(files, files2...)

	for _, file := range files {
		if file.IsDir() && file.Name() != "TEST12345" && file.Name() != ".locks" {
			url := getURLFromID(file.Name())
			urlList = append(urlList, url)
		}
	}

	return urlList
}

func main() {
	initConfig()

	conf := fetchConfig()

	fmt.Println("downloading using config.yml")

	path := getGamesPath(conf.ConfigYMLPath)
	urls := getGamesURLs(path)

	for index, url := range urls {
		printInfo("fetching URL %d: '%s'\n", index, url)

		// we need this because we can't verify the signature
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpClient := &http.Client{Transport: transport,
			Timeout: time.Duration(conf.DLTimeout) * time.Second}

		// TODO: retry logic goes here
		response, err := httpClient.Get(url)

		if isError(err) {
			printError("Error: Can't open url '%s'\n", url)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)

		if isError(err) {
			printError("can't read response body\n")
			break
		}
		patch := TitlePatch{}
		err = xml.Unmarshal([]byte(body), &patch)

		if isError(err) {
			printError("can't parse response XML\n")
			continue
		}

		printInfo("title '%s (%s) url '%s'\n",
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
