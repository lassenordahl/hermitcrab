package versions

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Manifest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Meta        map[string]interface{} `json:"metadata"`
	Versions    []Version              `json:"versions"`
}

func WriteManifestToFile(manifest Manifest, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(manifest)
	if err != nil {
		return err
	}

	return nil
}

func ParseManifestFromFile(filename string) (*Manifest, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m *Manifest) AddNewVersionFromFilePath(filePath string) ([]Version, error) {
	version, err := ParseFileName(filePath)
	if err != nil {
		return nil, err
	}

	return m.AddNewVersion(*version), nil
}

func (m *Manifest) AddNewVersion(version Version) []Version {
	m.Versions = append(m.Versions, version)
	return m.Versions
}

func (m *Manifest) DeduplicateVersions() []Version {
	seen := map[string]struct{}{}
	deduplicatedVersions := []Version{}

	for _, ver := range m.Versions {
		uniqueKey := fmt.Sprintf("%s-%d", ver.SemVersion, ver.BuildNumber)
		if _, ok := seen[uniqueKey]; !ok {
			deduplicatedVersions = append(deduplicatedVersions, ver)
		} else {
			continue
		}

		seen[uniqueKey] = struct{}{}
	}

	m.Versions = deduplicatedVersions
	return m.Versions
}

func (m *Manifest) SortVersions() []Version {
	versions := Versions(m.Versions)
	sort.Sort(versions)
	m.Versions = []Version(versions)
	return m.Versions
}

func (m *Manifest) GetLatestVersion() *Version {
	if len(m.Versions) == 0 {
		return nil
	}
	m.SortVersions()
	return &m.Versions[len(m.Versions)-1]
}

const skeletonManifest = "../testdata/base-no-versions.json"
const defaultManifest = "../testdata/base-versions.json"
const defaultOutput = "../testdata/output.json"

func AddNewVersionToManifest(
	prevManifestPath string, outputPath string, newVersionPath string,
) error {
	manifest, err := ParseManifestFromFile(prevManifestPath)
	if err != nil {
		return err
	}

	_, err = manifest.AddNewVersionFromFilePath(newVersionPath)
	if err != nil {
		return err
	}

	manifest.DeduplicateVersions()
	manifest.SortVersions()

	err = WriteManifestToFile(*manifest, outputPath)
	if err != nil {
		return err
	}

	return nil
}
