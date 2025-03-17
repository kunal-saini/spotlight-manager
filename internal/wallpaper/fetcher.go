package wallpaper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type SpotlightResponse struct {
	BatchRsp struct {
		Ver   string `json:"ver"`
		Items []struct {
			Item string `json:"item"`
		} `json:"items"`
	} `json:"batchrsp"`
}

type SpotlightItem struct {
	Ad struct {
		LandscapeImage struct {
			Asset string `json:"asset"`
		} `json:"landscapeImage"`
		PortraitImage struct {
			Asset string `json:"asset"`
		} `json:"portraitImage"`
		IconHoverText string `json:"iconHoverText"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		Copyright     string `json:"copyright"`
	} `json:"ad"`
}

type Fetcher struct {
	client *http.Client
	logger *log.Logger
}

func NewFetcher(logger *log.Logger) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (f *Fetcher) FetchSpotlightMetadata() (string, string, error) {
	// Microsoft Spotlight API endpoint
	url := "https://fd.api.iris.microsoft.com/v4/api/selection?&placement=88000820&bcnt=1&country=US&locale=en-US&fmt=json"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a typical browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var spotlightResp SpotlightResponse
	if err := json.Unmarshal(body, &spotlightResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(spotlightResp.BatchRsp.Items) == 0 {
		return "", "", fmt.Errorf("no items found in response")
	}

	var spotlightItem SpotlightItem
	if err := json.Unmarshal([]byte(spotlightResp.BatchRsp.Items[0].Item), &spotlightItem); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal item: %w", err)
	}

	imageURL := spotlightItem.Ad.LandscapeImage.Asset
	if imageURL == "" {
		return "", "", fmt.Errorf("no landscape image URL found")
	}

	metadata := fmt.Sprintf("%s\n%s\n%s",
		spotlightItem.Ad.Title,
		spotlightItem.Ad.Description,
		spotlightItem.Ad.Copyright)

	return imageURL, metadata, nil
}

func (f *Fetcher) DownloadImage(url string, destPath string) error {
	resp, err := f.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create the destination directory if it doesn't exist
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the destination file
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write image data: %w", err)
	}

	return nil
}
