package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all the configuration parameters for the application
type Config struct {
	QbittorrentUrl      string
	QbittorrentUsername string
	QbittorrentPassword string
	DeleteAfterMinutes  int
	PollIntervalSeconds int
}

// Reads environment variables and returns a Config instance
func New() (Config, error) {
	cfg := Config{
		QbittorrentUrl:      os.Getenv("QBITTORRENT_URL"),
		QbittorrentUsername: os.Getenv("QBITTORRENT_USERNAME"),
		QbittorrentPassword: os.Getenv("QBITTORRENT_PASSWORD"),
	}

	deleteAfterStr := os.Getenv("DELETE_AFTER_MINUTES")
	if deleteAfterStr == "" {
		return Config{}, fmt.Errorf("missing required environment variable: DELETE_AFTER_MINUTES")
	}
	deleteAfter, err := strconv.Atoi(deleteAfterStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DELETE_AFTER_MINUTES value: %w", err)
	}
	cfg.DeleteAfterMinutes = deleteAfter

	pollIntervalStr := os.Getenv("POLL_INTERVAL_SECONDS")
	if pollIntervalStr == "" {
		return Config{}, fmt.Errorf("missing required environment variable: POLL_INTERVAL_SECONDS")
	}
	pollInterval, err := strconv.Atoi(pollIntervalStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid POLL_INTERVAL_SECONDS value: %w", err)
	}
	cfg.PollIntervalSeconds = pollInterval

	// Validate required fields
	missing := []string{}
	if cfg.QbittorrentUrl == "" {
		missing = append(missing, "QBITTORRENT_URL")
	}
	if cfg.QbittorrentUsername == "" {
		missing = append(missing, "QBITTORRENT_USERNAME")
	}
	if cfg.QbittorrentPassword == "" {
		missing = append(missing, "QBITTORRENT_PASSWORD")
	}

	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}
