package profiles

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"slv.sh/slv/internal/core/config"
)

const (
	configHTTPUrlKey    = "url"
	configHTTPHeaderKey = "auth-header"
)

var httpArgs = []arg{
	{
		name:        configHTTPUrlKey,
		required:    true,
		description: "The HTTP base URL of the remote profile",
	},
	{
		name:        configHTTPHeaderKey,
		sensitive:   true,
		description: "The header to be used for HTTP URLs protected by auth. E.g. 'Authorization: Bearer <token>'",
	},
}

func getUrlContent(url, header string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", config.AppNameUpperCase+"-"+config.Version+" ("+runtime.GOOS+"/"+runtime.GOARCH+")")
	if header != "" {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format: %s", header)
		}
		req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	return bodyBytes, nil
}

func httpPull(dir string, config map[string]string) (err error) {
	url := strings.TrimSuffix(config[configHTTPUrlKey], "/")
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("invalid HTTP URL: %s", url)
	}
	header := config[configHTTPHeaderKey]
	var envManifestFileName, settingsFileName string
	var environmentBytes, settingsBytes []byte
	for _, envManifestFileName = range envManifestFileNames {
		if environmentBytes, err = getUrlContent(url+"/"+envManifestFileName, header); err != nil {
			return fmt.Errorf("failed to get environments manifest: %w", err)
		} else if len(environmentBytes) > 0 {
			break
		}
	}
	for _, settingsFileName = range settingsFileNames {
		if settingsBytes, err = getUrlContent(url+"/"+settingsFileName, header); err != nil {
			return fmt.Errorf("failed to get settings manifest: %w", err)
		} else if len(settingsBytes) > 0 {
			break
		}
	}
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}
	if len(environmentBytes) > 0 {
		if err = os.WriteFile(dir+"/"+envManifestFileName, environmentBytes, 0644); err != nil {
			return fmt.Errorf("failed to write environments manifest: %w", err)
		}
	}
	if len(settingsBytes) > 0 {
		if err = os.WriteFile(dir+"/"+settingsFileName, settingsBytes, 0644); err != nil {
			return fmt.Errorf("failed to write settings manifest: %w", err)
		}
	}
	return nil
}

func httpSetup(dir string, config map[string]string) (err error) {
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create profile data directory: %w", err)
	}
	if err = httpPull(dir, config); err != nil {
		os.RemoveAll(dir)
	}
	return
}
