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
	"github.com/mattn/go-zglob"
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
	printInfo("Parsing '" + configYML)
	path := ""
	file, err := os.Open(configYML)

	if isError(err) {
		printError(fmt.Sprintf("Couldn't open '%s' (errorcode: %s)\n", configYML, err))
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
			printDebug("emudir: " + emulatorDir)
		}
		if strings.Contains(line, "/dev_hdd0/") {
			path = strings.Replace(strings.TrimSpace(strings.Split(line, ":")[1]), "$(EmulatorDir)", emulatorDir, -1)
			printDebug("path: " + path)
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

/* gets the game's version */

func getVersion(path string) string {
	// finds the PARAM.SFO
	params, err := zglob.Glob(path + "/**/PARAM.SFO")
	if isError(err) {
		printError("Couldn't find "+path+"*PARAM.sfo  (errorcode: %s)\n", err)
		return ""
	}
	param := params[0]
	file, err := os.Open(param)
	defer file.Close()

	if isError(err) {
		printError(fmt.Sprintf("Couldn't open '%s' (errorcode: %s)\n", param, err))
		return ""
	}
	// goes to 16 bytes before the end
	file.Seek(-8, 2)
	buf := make([]byte, 6)
	file.Read(buf)
	version := string(buf[:5])
	printDebug("The version for '%s' is : %s", path, version)

	return version
}

/* gets games URLs and versions from a specific folder */

func getGamesFromFolder(path string) []GameInfo {
	var games []GameInfo
	files, err := ioutil.ReadDir(path)
	if err != nil {
		printError(fmt.Sprintf("Couldn't open '%s' (errorcode: '%s')\n", path, err))
		return games
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != "TEST12345" && file.Name() != ".locks" {
			url := getURLFromID(file.Name())
			version := getVersion(path + file.Name())
			game := GameInfo{
				ID:      file.Name(),
				URL:     url,
				Version: version,
			}
			games = append(games, game)
		}
	}
	return games
}

/* gets games URLs and versions from the various folders */

func getGames(path string) []GameInfo {
	// first from the disc folder
	games := getGamesFromFolder(path + "disc/")

	// then reads the game folder
	games = append(games, getGamesFromFolder(path+"game/")...)

	return games
}

func main() {
	initConfig()

	conf := fetchConfig()

	printInfo("downloading using config.yml")

	path := getGamesPath(conf.ConfigYMLPath)
	games := getGames(path)


	for index, game := range games {
		url := game.URL
		printInfo("fetching URL %d: '%s'", index, url)

		// we need this because we can't verify the signature
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpClient := &http.Client{Transport: transport,
			Timeout: time.Duration(conf.DLTimeout) * time.Second}

		// TODO: retry logic goes here
		response, err := httpClient.Get(url)

		if isError(err) {
			printError("Error: Can't open url '%s'", url)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)

		if isError(err) {
			printError("can't read response body.")
			break
		}
		patch := TitlePatch{}
		err = xml.Unmarshal([]byte(body), &patch)

		if isError(err) {
			printError("can't parse response XML.")
			continue
		}

		printInfo("title '%s (%s) url '%s'",
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
