package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dorik33/Test/internal/models"
)

func FetchSongInfo(BaseURL string, req models.SongRequest) (*models.SongDetail, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/info?group=%s&song=%s",
		BaseURL,
		url.QueryEscape(req.GroupName),
		url.QueryEscape(req.SongName),
	)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &songDetail, nil
}
