package versions

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFileName(t *testing.T) {
	type args struct {
		fn string
	}
	type expected struct {
		version     string
		filename    string
		prerelease  string
		buildNumber int
	}

	tests := []struct {
		name    string
		args    args
		want    expected
		wantErr error
	}{
		{
			name: "file is the semver with a v prefix",
			args: args{
				fn: "v1.2.34",
			},
			want: expected{
				version:  "1.2.34",
				filename: "v1.2.34",
			},
			wantErr: nil,
		},
		{
			name: "file is the semver without a v prefix",
			args: args{
				fn: "24.3.1",
			},
			want: expected{
				version:  "24.3.1",
				filename: "24.3.1",
			},
			wantErr: nil,
		},
		{
			name: "dbconsole tar zst naming convention",
			args: args{
				fn: "ui_v24.3.1.tar.zst",
			},
			want: expected{
				version:  "24.3.1",
				filename: "ui_v24.3.1.tar.zst",
			},
			wantErr: nil,
		},
		{
			name: "molt ui naming convention",
			args: args{
				fn: "molt-0.1.1.tar.gz",
			},
			want: expected{
				version:  "0.1.1",
				filename: "molt-0.1.1.tar.gz",
			},
			wantErr: nil,
		},
		{
			name: "dbconsole file naming convention with build number",
			args: args{
				fn: "24.1.5-ui.1.tar.gz",
			},
			want: expected{
				version:    "24.1.5-ui.1",
				filename:   "24.1.5-ui.1.tar.gz",
				prerelease: "ui.1",
			},
			wantErr: nil,
		},
		{
			name: "dbconsole file naming convention with prerelease and number",
			args: args{
				fn: "24.1.5-beta+100.tar.gz",
			},
			want: expected{
				version:     "24.1.5-beta+100",
				filename:    "24.1.5-beta+100.tar.gz",
				prerelease:  "beta",
				buildNumber: 100,
			},
			wantErr: nil,
		},
		{
			name: "full file path and pre-release and build number",
			args: args{
				fn: "./artifacts/test/24.1.5-beta+100.tar.gz",
			},
			want: expected{
				version:     "24.1.5-beta+100",
				filename:    "./artifacts/test/24.1.5-beta+100.tar.gz",
				prerelease:  "beta",
				buildNumber: 100,
			},
			wantErr: nil,
		},
		{
			name: "full website path",
			args: args{
				fn: "https://server.com/best/24.1.5-beta+100.tar.gz",
			},
			want: expected{
				version:     "24.1.5-beta+100",
				filename:    "https://server.com/best/24.1.5-beta+100.tar.gz",
				prerelease:  "beta",
				buildNumber: 100,
			},
			wantErr: nil,
		},
		{
			name: "full website path",
			args: args{
				fn: "https://server.com/best/24.1.5-beta+100.tar.gz",
			},
			want: expected{
				version:     "24.1.5-beta+100",
				filename:    "https://server.com/best/24.1.5-beta+100.tar.gz",
				prerelease:  "beta",
				buildNumber: 100,
			},
			wantErr: nil,
		},
		{
			name: "invalid build number",
			args: args{
				fn: "https://server.com/best/24.1.5-beta+new.tar.gz",
			},
			want:    expected{},
			wantErr: errors.New(`strconv.Atoi: parsing "new": invalid syntax`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFileName(tt.args.fn)
			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.filename, got.FilePath)
				require.Equal(t, tt.want.version, got.SemVersion)
				require.Equal(t, tt.want.prerelease, got.Prerelease())
				require.Equal(t, tt.want.buildNumber, got.BuildNumber)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid version",
			args: args{
				version: "24.1.5-beta+100",
			},
			wantErr: nil,
		},
		{
			name: "valid version path",
			args: args{
				version: "https://gcs.com/molt-24.1.5-ui+1",
			},
			wantErr: nil,
		},
		{
			name: "invalid version",
			args: args{
				version: "invalid--version",
			},
			wantErr: errors.New("no semantic version found in filename"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.args.version)
			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
