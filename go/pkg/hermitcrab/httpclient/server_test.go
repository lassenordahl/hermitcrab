package httpclient

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab/bucket"
	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeLatestVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBM := bucket.NewMockBucketManager(ctrl)
	server := NewTestServer(t, WithBucketManager(mockBM))

	latestVersion, err := version.ParseVersion("24.1.0-ui.1")
	require.NoError(t, err)

	mockBM.EXPECT().GetLatestPatchVersion(gomock.Any(), "24.1").Return(latestVersion, nil)
	mockBM.EXPECT().DownloadPatchVersion(gomock.Any(), latestVersion).Return(createMockTarGz("index.html", "<html>Latest Version</html>"), nil)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Latest Version")
}

func TestServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBM := bucket.NewMockBucketManager(ctrl)
	server := NewTestServer(t, WithBucketManager(mockBM))

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		setupMock      func()
	}{
		{
			name:           "Latest Version",
			path:           "/",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				latestVersion, _ := version.ParseVersion("24.1.0-ui.1")
				mockBM.EXPECT().GetLatestPatchVersion(gomock.Any(), "24.1").Return(latestVersion, nil)
				mockBM.EXPECT().DownloadPatchVersion(gomock.Any(), latestVersion).Return(createMockTarGz("index.html", "<html>Latest Version</html>"), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req, _ := http.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

// ... (createMockTarGz function remains the same)
func createMockTarGz(filename, content string) io.ReadCloser {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(content)),
	}
	tw.WriteHeader(hdr)
	tw.Write([]byte(content))

	tw.Close()
	gw.Close()

	return io.NopCloser(bytes.NewReader(buf.Bytes()))
}
