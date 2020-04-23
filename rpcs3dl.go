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
)

// just a generic helper

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

/* parses the given config.yml file and returns the path to dev_hdd0 */

func getGamesPath(configYML string) string {
	fmt.Printf("Parsing '%s'\n", configYML)
	path := "test"
	file, err := os.Open(configYML)

	if isError(err) {
		fmt.Printf("Couldn't open '%s' (errorcode: %d)", configYML, err)
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
			fmt.Printf("emudir: " + emulatorDir + "TT\n")
		}
		if strings.Contains(line, "/dev_hdd0/") {
			path = strings.Replace(strings.TrimSpace(strings.Split(line, ":")[1]), "$(EmulatorDir)", emulatorDir, -1)
			fmt.Printf("path: " + path + "\n")
		}

		// we use err to figure out end of input
		if isError(err) {
			return path
		}
	}

	return path
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
		fmt.Printf("Couldn't open '%s' (errorcode: %d)", path, err)
		return urlList
	}
	// then reads the game folder
	files2, err := ioutil.ReadDir(path + "game")
	if err != nil {
		fmt.Printf("Couldn't open '%s' (errorcode: %d)", path, err)
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
	fmt.Println("downloading using config.yml")

	path := getGamesPath(conf.ConfigYMLPath)
	urls := getGamesURLs(path)

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
			continue
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
