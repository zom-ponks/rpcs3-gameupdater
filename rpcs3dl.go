/* rpcs3-gameupdater main entry point */

package main

import (
	// TODO: these are the UI libs
	//"fyne.io/fyne/app"
	//"fyne.io/fyne/widget"
	// TODO: figure out if we really need this

	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

/* parses the given config.yml file and returns the path to dev_hdd0 */

func getLocalGamesPath(configYML string) string {
	printInfo("Parsing '" + configYML)
	path := ""
	file, err := os.Open(configYML)

	if err != nil {
		printError("Couldn't open '%s' (errorcode: %s)", configYML, err)
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
		if err != nil {
			return path
		}
	}
}

/* replaces the variable in the URL with the gameID */

func getURLFromID(id string) string {
	return fmt.Sprintf(urlPattern, id, id)
}

/* gets the game's version */

func getCategoryAndVersion(path string) (string, float64) {
	var folder string
	if strings.Contains(path, "/disc/") {
		folder = "/PS3_GAME"
	}
	// finds the PARAM.SFO
	params, err := filepath.Glob(path + folder + "/PARAM.SFO")
	if err != nil {
		printError("Error finding %s/**/PARAM.sfo  (errorcode: %s)", path, err)
		return "", 0.0
	}
	param := params[0]
	file, err := os.Open(param)
	defer file.Close()

	if err != nil {
		printError("Couldn't open '%s' (errorcode: %s)", param, err)
		return "", 0.0
	}
	kvp := readParamSFO(file)
	cat := getCategory(kvp)
	ver := getAppVersion(kvp)
	// in case there is no app version, use version instead
	if ver == "" {
		ver = getVersion(kvp)
	}
	versionF, err := strconv.ParseFloat(ver[0:5], 64)
	if err != nil {
		printError("Couldn't convert '%s' (errorcode: '%s')", ver, err)
	}
	return cat, versionF
}

/* gets games URLs and versions from a specific folder */

func getLocalGamesFromFolder(games map[string]*GameInfo, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		printError("Couldn't open '%s' (errorcode: '%s')", path, err)
		return
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != "TEST12345" && file.Name() != ".locks" && !strings.Contains(file.Name(), "INST") {
			url := getURLFromID(file.Name())

			category, version := getCategoryAndVersion(path + file.Name())

			if game, ok := games[file.Name()]; ok {
				if game.Version < version {
					game.Version = version
				}
			} else {
				game := GameInfo{
					Category: category,
					URL:      url,
					Version:  version,
				}
				games[file.Name()] = &game
			}
		}
	}
}

/* gets games URLs and versions from the various folders */

func getLocalGames(path string) map[string]*GameInfo {
	// first from the disc folder
	games := make(map[string]*GameInfo)
	getLocalGamesFromFolder(games, path+"disc/")

	// then the game folder
	getLocalGamesFromFolder(games, path+"game/")

	return games
}

func getGamesFromServer(games map[string]*GameInfo) {
	var downloaded []string
	count := 0
	var wg sync.WaitGroup
	wg.Add(len(games))
	for gameID, game := range games {
		count = count + 1
		go func(count int, gameID string, game *GameInfo) {
			defer wg.Done()
			printDebug("gameID: %s, url: %s, version: %f", gameID, game.URL, game.Version)
			printInfo("Downloading PKGs for game %d/%d, at URL: '%s'", count, len(games), game.URL)

			patch := TitlePatch{}
			err := xml.Unmarshal(getXML(game.URL), &patch)

			if err != nil {
				printError("can't parse response XML.")
				return
			}

			count2 := 1
			for i := range patch.Tag.Package {
				printInfo("Downloading PKG %d/%d, for game %s", count2, len(patch.Tag.Package), gameID)
				count2 = count2 + 1
				printDebug("title '%s' (%s) version %s url '%s' SHA '%s':",
					patch.Tag.Package[i].Paramsfo.TITLE,
					patch.Titleid,
					patch.Tag.Package[i].Version,
					patch.Tag.Package[i].URL,
					patch.Tag.Package[i].SHA1)
				version, err := strconv.ParseFloat(patch.Tag.Package[i].Version, 64)
				if err != nil {
					printError("Couldn't convert '%s' (errorcode: '%s')", patch.Tag.Package[i].Version, err)
				}
				if version < game.Version {
					printDebug("Version %f is inferior to current of %f", version, game.Version)
					return
				}
				if getPKG(patch.Tag.Package[i].URL, patch.Tag.Package[i].SHA1) {
					downloaded = append(downloaded, gameID)
				}
			}
		}(count, gameID, game)
	}
	wg.Wait()
	printInfo("We've downloaded %d games", len(downloaded))
}

func main() {
	parseArguments()
	initConfig()
	initDownloader()
	path := getLocalGamesPath(fetchConfig().ConfigYMLPath)
	games := getLocalGames(path)
	getGamesFromServer(games)

	// test
	fmt.Printf("Terminal: %v\n", isTTY())

	// TODO: UI stuff
	//app := app.New()

	//w := app.NewWindow("Hello")
	//w.SetContent(widget.NewLabel("Hello Fyne!"))

	//w.ShowAndRun()
}
