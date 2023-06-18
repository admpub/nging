package osinfo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Release contains all the info about os in nested structs
// All values in the top level of the struct are guaranteed
// in the sense that they are always returned.
// Nested structs are only populated for the os they are compiled for
// IE: Release.Windows will only have values if you are on windows
type Release struct {
	Runtime string
	Arch    string
	Name    string
	Version string
	Windows windowsRelease
	Linux   linuxRelease
	BSD     bsdRelease
	MacOs   macOsRelease
}

type windowsRelease struct {
	Build string
}

type linuxRelease struct {
	Distro     string
	Kernel     string
	PkgManager string
}

type bsdRelease struct {
	Kernel     string
	PkgManager string
}

type macOsRelease struct {
	VersionName string
}

// Provides a preformated print for every support OS
func (r *Release) String() string {
	b := bytes.NewBuffer(nil)

	fmt.Fprintln(b, "Runtime:", r.Runtime)
	fmt.Fprintln(b, "Architecture:", r.Arch)
	fmt.Fprintln(b, "OS Name:", r.Name)
	fmt.Fprintln(b, "Version:", r.Version)

	switch r.Runtime {
	case "windows":
		fmt.Fprintln(b, "Build Number:", r.Windows.Build)
	case "freebsd":
		fmt.Fprintln(b, "Kernel:", r.BSD.Kernel)
	case "openbsd":
		fmt.Fprintln(b, "Kernel:", r.BSD.Kernel)
	case "linux":
		fmt.Fprintln(b, "Kernel:", r.Linux.Kernel)
		fmt.Fprintln(b, "Distro:", r.Linux.Distro)
		fmt.Fprintln(b, "Package Manager:", r.Linux.PkgManager)
	}

	return b.String()
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, `"`, ``)

	return strings.TrimSpace(s)
}

func pathExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func readFile(filename string) string {
	content, _ := ioutil.ReadFile(filename)
	return string(content)
}
