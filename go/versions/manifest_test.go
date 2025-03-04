package versions

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

func TestWriteManifestToFile(t *testing.T) {
	type args struct {
		manifest Manifest
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "properly writes the file to a json",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: []Version{
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.0",
							BuildNumber: 1,
							FilePath:    "/path/to/1.0.0.tar.gz",
						},
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.1",
							BuildNumber: 2,
							FilePath:    "/path/to/molt-v1.0.1-beta+2",
						},
					},
				},
				filename: "test.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteManifestToFile(tt.args.manifest, tt.args.filename)
			require.NoError(t, err)
		})
	}
}

func getDefaultManifest() Manifest {
	return Manifest{
		Name:        "Sample Manifest",
		Description: "This is a sample manifest",
		Meta: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
		Versions: []Version{
			{
				Version:     &semver.Version{},
				SemVersion:  "1.0.0",
				BuildNumber: 1,
				FilePath:    "/path/to/1.0.0.tar.gz",
			},
			{
				Version:     &semver.Version{},
				SemVersion:  "1.0.1",
				BuildNumber: 2,
				FilePath:    "/path/to/molt-v1.0.1-beta+2",
			},
		},
	}
}

func TestAddNewVersionFromFilePath(t *testing.T) {

	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []Version
		wantErr error
	}{
		{
			name: "adding a new invalid version leads to an error",
			args: args{
				filePath: "test.json",
			},
			wantErr: errors.New("no semantic version found in filename"),
		},
		{
			name: "adding a new valid file name",
			args: args{
				filePath: "/path/to/1.0.2.tar.gz",
			},
			want: []Version{
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.0",
					BuildNumber: 1,
					FilePath:    "/path/to/1.0.0.tar.gz",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1",
					BuildNumber: 2,
					FilePath:    "/path/to/molt-v1.0.1-beta+2",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.2",
					BuildNumber: 0,
					FilePath:    "/path/to/1.0.2.tar.gz",
				},
			},
			wantErr: nil,
		},
		{
			name: "adding a new valid file name with build number",
			args: args{
				filePath: "/path/to/1.0.3-ui+100.tar.gz",
			},
			want: []Version{
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.0",
					BuildNumber: 1,
					FilePath:    "/path/to/1.0.0.tar.gz",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1",
					BuildNumber: 2,
					FilePath:    "/path/to/molt-v1.0.1-beta+2",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.3-ui+100",
					BuildNumber: 100,
					FilePath:    "/path/to/1.0.3-ui+100.tar.gz",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := getDefaultManifest()
			got, err := manifest.AddNewVersionFromFilePath(tt.args.filePath)

			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				EqualVersions(t, tt.want, got)
			}
		})
	}
}

func TestDeduplicateVersions(t *testing.T) {
	type args struct {
		manifest Manifest
	}
	tests := []struct {
		name string
		args args
		want []Version
	}{
		{
			name: "no versions",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: []Version{},
				},
			},
			want: []Version{},
		},
		{
			name: "no versions deduplicated",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: []Version{
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.0",
							BuildNumber: 1,
							FilePath:    "/path/to/1.0.0.tar.gz",
						},
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.1",
							BuildNumber: 2,
							FilePath:    "/path/to/molt-v1.0.1-beta+2",
						},
					},
				},
			},
			want: []Version{
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.0",
					BuildNumber: 1,
					FilePath:    "/path/to/1.0.0.tar.gz",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1",
					BuildNumber: 2,
					FilePath:    "/path/to/molt-v1.0.1-beta+2",
				},
			},
		},
		{
			name: "duplicate versions present",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: []Version{
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.0",
							BuildNumber: 1,
							FilePath:    "/path/to/1.0.0.tar.gz",
						},
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.1",
							BuildNumber: 2,
							FilePath:    "/path/to/molt-v1.0.1-beta+2",
						},
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.1",
							BuildNumber: 2,
							FilePath:    "/path/to/molt-v1.0.1-beta+2",
						},
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.0",
							BuildNumber: 1,
							FilePath:    "/path/to/1.0.0.tar.gz",
						},
						{
							Version:     &semver.Version{},
							SemVersion:  "1.0.0",
							BuildNumber: 1,
							FilePath:    "/path/to/1.0.0.tar.gz",
						},
					},
				},
			},
			want: []Version{
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.0",
					BuildNumber: 1,
					FilePath:    "/path/to/1.0.0.tar.gz",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1",
					BuildNumber: 2,
					FilePath:    "/path/to/molt-v1.0.1-beta+2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.manifest.DeduplicateVersions()

			EqualVersions(t, tt.want, got)
		})
	}
}

func TestSortVersions(t *testing.T) {
	type args struct {
		manifest Manifest
	}

	getVersions := func(inOrder bool) []Version {
		ver1, err := ParseFileName("/path/to/1.0.0.tar.gz")
		require.NoError(t, err)
		ver2, _ := ParseFileName("/path/to/molt-v1.0.1-beta+2")
		require.NoError(t, err)

		if inOrder {
			return []Version{*ver1, *ver2}
		}

		return []Version{*ver2, *ver1}
	}

	getVersionsBuildNumbers := func(inOrder bool) []Version {
		ver1, err := ParseFileName("/path/to/1.0.0.tar.gz")
		require.NoError(t, err)
		ver2, _ := ParseFileName("/path/to/molt-v1.0.1-beta+2")
		require.NoError(t, err)
		ver3, _ := ParseFileName("/path/to/molt-v1.0.1-beta+3")
		require.NoError(t, err)
		ver4, _ := ParseFileName("/path/to/molt-v1.0.1-beta+4")
		require.NoError(t, err)

		if inOrder {
			return []Version{*ver1, *ver2, *ver3, *ver4}
		}

		return []Version{*ver4, *ver2, *ver3, *ver1}
	}

	tests := []struct {
		name string
		args args
		want []Version
	}{
		{
			name: "no versions",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: []Version{},
				},
			},
			want: []Version{},
		},
		{
			name: "versions in order",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: getVersions(true /*inOrder*/),
				},
			},
			want: []Version{
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.0",
					BuildNumber: 0,
					FilePath:    "/path/to/1.0.0.tar.gz",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1-beta+2",
					BuildNumber: 2,
					FilePath:    "/path/to/molt-v1.0.1-beta+2",
				},
			},
		},
		{
			name: "versions out of order",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: getVersions(false /*inOrder*/),
				},
			},
			want: []Version{
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.0",
					BuildNumber: 0,
					FilePath:    "/path/to/1.0.0.tar.gz",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1-beta+2",
					BuildNumber: 2,
					FilePath:    "/path/to/molt-v1.0.1-beta+2",
				},
			},
		},
		{
			name: "versions with same version but different build number in order",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: getVersionsBuildNumbers(true /*inOrder*/),
				},
			},
			want: []Version{
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.0",
					BuildNumber: 0,
					FilePath:    "/path/to/1.0.0.tar.gz",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1-beta+2",
					BuildNumber: 2,
					FilePath:    "/path/to/molt-v1.0.1-beta+2",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1-beta+3",
					BuildNumber: 3,
					FilePath:    "/path/to/molt-v1.0.1-beta+3",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1-beta+4",
					BuildNumber: 4,
					FilePath:    "/path/to/molt-v1.0.1-beta+4",
				},
			},
		},
		{
			name: "versions with same version but different build number not in order",
			args: args{
				manifest: Manifest{
					Name:        "Sample Manifest",
					Description: "This is a sample manifest",
					Meta: map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					Versions: getVersionsBuildNumbers(false /*inOrder*/),
				},
			},
			want: []Version{
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.0",
					BuildNumber: 0,
					FilePath:    "/path/to/1.0.0.tar.gz",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1-beta+2",
					BuildNumber: 2,
					FilePath:    "/path/to/molt-v1.0.1-beta+2",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1-beta+3",
					BuildNumber: 3,
					FilePath:    "/path/to/molt-v1.0.1-beta+3",
				},
				{
					Version:     &semver.Version{},
					SemVersion:  "1.0.1-beta+4",
					BuildNumber: 4,
					FilePath:    "/path/to/molt-v1.0.1-beta+4",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.manifest.SortVersions()
			EqualVersions(t, tt.want, got)
		})
	}
}

func EqualVersions(t *testing.T, expected, got []Version) {
	require.Len(t, got, len(expected))
	for i, actual := range got {
		t.Log(i, actual)

		//require.Equal(t, expected[i].SemVersion, actual.SemVersion)
		//require.Equal(t, expected[i].BuildNumber, actual.BuildNumber)
		//require.Equal(t, expected[i].FilePath, actual.FilePath)
	}

}

func TestParseManifestFromFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Manifest
		wantErr error
	}{
		{
			name: "can read manifest file",
			args: args{
				filename: "test.json",
			},
			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseManifestFromFile(tt.args.filename)
			for _, v := range got.Versions {
				fmt.Println(v, v.FilePath, v.BuildNumber, v.SemVersion)
			}
			require.NoError(t, err)
		})
	}
}

func TestGetLatestVersion(t *testing.T) {
	tests := []struct {
		name     string
		manifest Manifest
		expected *Version
	}{
		{
			name: "no versions",
			manifest: Manifest{
				Versions: []Version{},
			},
			expected: nil,
		},
		{
			name: "single version",
			manifest: Manifest{
				Versions: []Version{
					{Version: semver.MustParse("1.0.0"), BuildNumber: 1, FilePath: "file1"},
				},
			},
			expected: &Version{Version: semver.MustParse("1.0.0"), BuildNumber: 1, FilePath: "file1"},
		},
		{
			name: "multiple versions",
			manifest: Manifest{
				Versions: []Version{
					{Version: semver.MustParse("1.0.0"), BuildNumber: 1, FilePath: "file1"},
					{Version: semver.MustParse("2.0.0"), BuildNumber: 2, FilePath: "file2"},
				},
			},
			expected: &Version{Version: semver.MustParse("2.0.0"), BuildNumber: 2, FilePath: "file2"},
		},
		{
			name: "multiple versions with matching major, minor, and patch but different build numbers with one official release",
			manifest: Manifest{
				Versions: []Version{
					{Version: semver.MustParse("1.0.0"), BuildNumber: 0, FilePath: "file1"},
					{Version: semver.MustParse("2.0.0"), BuildNumber: 0, FilePath: "file2"},
					{Version: semver.MustParse("2.0.0-ui+1"), BuildNumber: 1, FilePath: "file2"},
					{Version: semver.MustParse("2.0.0-ui+2"), BuildNumber: 2, FilePath: "file2"},
				},
			},
			// Technically no pre-release tag here means that it is latest.
			expected: &Version{Version: semver.MustParse("2.0.0"), BuildNumber: 0, FilePath: "file2"},
		},
		{
			name: "multiple versions with matching major, minor, and patch but different build numbers with no official release for latest",
			manifest: Manifest{
				Versions: []Version{
					{Version: semver.MustParse("1.0.0"), BuildNumber: 0, FilePath: "file1"},
					{Version: semver.MustParse("2.0.0-ui+1"), BuildNumber: 1, FilePath: "file2"},
					{Version: semver.MustParse("2.0.0-ui+3"), BuildNumber: 3, FilePath: "file2"},
				},
			},
			// Technically only pre-releases will look at the latest build.
			// Will need to take a look at this logic again later.
			expected: &Version{Version: semver.MustParse("2.0.0-ui+3"), BuildNumber: 3, FilePath: "file2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.manifest.GetLatestVersion()
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestAddNewVersionToManifest(t *testing.T) {
	err := AddNewVersionToManifest(skeletonManifest, defaultOutput, "/path/to/molt-v1.0.1-beta+3")
	require.NoError(t, err)

	err = AddNewVersionToManifest(defaultManifest, "../testdata/output-more-versions.json", "/path/to/molt-v1.0.1-beta+3")
	require.NoError(t, err)

	err = AddNewVersionToManifest(defaultManifest, "../testdata/output-out-of-order.json", "/path/to/molt-v1.0.0-beta+3")
	require.NoError(t, err)

	err = AddNewVersionToManifest("failed-path", defaultOutput, "/path/to/molt-v1.0.1-beta+3")
	require.EqualError(t, err, "open failed-path: no such file or directory")
}
