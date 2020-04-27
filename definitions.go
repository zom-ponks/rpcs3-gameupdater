/* rpcs3-gameupdater - constants & data structures */

package main

/* constants */

const urlPattern string = "https://a0.ww.np.dl.playstation.net/tpl/np/%s/%s-ver.xml"
const gamesYAML = "games.yml"
const userAgent string = "Mozilla/5.0 (PLAYSTATION 3; 3.55)"

/* structures */

type GameInfo struct {
	ID      string
	URL     string
	Version string
}

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
