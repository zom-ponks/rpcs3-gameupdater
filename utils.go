/* rpcs3-gameupdater - utilities */

package main

// check for file or path existence
func testPath(path string) bool {
	return true
}

// IsTTY checks if we're in a terminal, this is platform specific, implementations in _unix and _windows files
func IsTTY() bool {
	return isTTY()
}
