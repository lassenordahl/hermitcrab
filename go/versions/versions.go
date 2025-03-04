package versions

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	semver "github.com/Masterminds/semver/v3"
)

type Versions []Version

type Version struct {
	*semver.Version `json:"-"`
	SemVersion      string `json:"version"`
	BuildNumber     int    `json:"build_number"`
	FilePath        string `json:"filename"`
}

// Implement sort interface
// Len returns the length of a collection. The number of Version instances
// on the slice.
func (v Versions) Len() int {
	return len(v)
}

// Less is needed for the sort interface to compare two Version objects on the
// slice. If checks if one is less than the other.
func (v Versions) Less(i, j int) bool {
	if v[i].Version.Equal(v[j].Version) {
		return v[i].BuildNumber < v[j].BuildNumber
	}
	return v[i].Version.LessThan(v[j].Version)
}

// Swap is needed for the sort interface to replace the Version objects
// at two different positions in the slice.
func (c Versions) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (v *Version) MarshalJSON() ([]byte, error) {
	type Alias struct {
		SemVersion  string `json:"version"`
		BuildNumber int    `json:"build_number"`
		FilePath    string `json:"filename"`
	}

	aux := &Alias{
		SemVersion:  v.SemVersion,
		BuildNumber: v.BuildNumber,
		FilePath:    v.FilePath,
	}

	return json.Marshal(aux)
}

func (v *Version) UnmarshalJSON(data []byte) error {
	type Alias struct {
		SemVersion  string `json:"version"`
		BuildNumber int    `json:"build_number"`
		FilePath    string `json:"filename"`
	}

	aux := &Alias{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	v.SemVersion = aux.SemVersion
	v.BuildNumber = aux.BuildNumber
	v.FilePath = aux.FilePath
	ver, err := ParseFileName(v.FilePath)
	if err != nil {
		return err
	}
	v.Version = ver.Version

	return nil
}

var trimPatterns = []string{"tar.gz", "tar.zst"}

const semverPattern string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

var semverRegex = regexp.MustCompile(semverPattern)

func parseSemanticVersion(filename string) (string, error) {
	// Define the semantic versioning regex pattern

	// Find the first match in the filename
	match := semverRegex.FindString(filename)
	if match == "" {
		return "", fmt.Errorf("no semantic version found in filename")
	}

	return match, nil
}

// ParseString parses a version filename into a Version struct.
func ParseFileName(filePath string) (*Version, error) {
	// Only get the base, which is the file name.
	fileName := path.Base(filePath)
	semVer, err := parseSemanticVersion(fileName)
	if err != nil {
		return nil, err
	}

	// Trim the file extensions, if relevant.
	// We want to remove the file information because
	// in can mess up how the version presents as.
	for _, pattern := range trimPatterns {
		semVer = strings.TrimRight(semVer, pattern)
	}

	v, err := semver.NewVersion(semVer)
	if err != nil {
		return nil, err
	}

	buildNum := 0
	if v.Metadata() != "" {
		parsedNum, err := strconv.Atoi(v.Metadata())
		if err != nil {
			return nil, err
		}
		buildNum = parsedNum
	}

	return &Version{
		Version:     v,
		SemVersion:  v.String(),
		FilePath:    filePath,
		BuildNumber: buildNum,
	}, err
}

func ValidateVersion(version string) error {
	fileName := path.Base(version)
	semVer, err := parseSemanticVersion(fileName)
	if err != nil {
		return err
	}

	// Trim the file extensions, if relevant.
	// We want to remove the file information because
	// in can mess up how the version presents as.
	for _, pattern := range trimPatterns {
		semVer = strings.TrimRight(semVer, pattern)
	}

	ver, err := semver.NewVersion(semVer)
	if err != nil {
		return err
	}

	fmt.Printf("validated version: %s\n", ver.String())
	return nil
}
