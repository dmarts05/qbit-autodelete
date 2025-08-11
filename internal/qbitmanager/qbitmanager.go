package qbitmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/dmarts05/qbit-autodelete/internal/config"
)

// QbitManager manages qBittorrent connection and torrent operations
type QbitManager struct {
	cfg             config.Config
	client          *http.Client
	cookies         []*http.Cookie
	completionTimes map[string]time.Time
}

// Creates a new QbitManager from config
func New(cfg config.Config) (*QbitManager, error) {
	m := &QbitManager{
		cfg:             cfg,
		client:          &http.Client{},
		completionTimes: make(map[string]time.Time),
	}

	if err := m.login(); err != nil {
		return nil, err
	}
	return m, nil
}

// Authenticates with qBittorrent
func (m *QbitManager) login() error {
	data := url.Values{}
	data.Set("username", m.cfg.QbittorrentUsername)
	data.Set("password", m.cfg.QbittorrentPassword)

	resp, err := m.client.PostForm(m.cfg.QbittorrentUrl+"/api/v2/auth/login", data)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed, status code: %d", resp.StatusCode)
	}

	m.cookies = resp.Cookies()
	return nil
}

// Gets the list of torrents from qBittorrent
func (m *QbitManager) getTorrents() ([]Torrent, error) {
	req, _ := http.NewRequest("GET", m.cfg.QbittorrentUrl+"/api/v2/torrents/info", nil)
	for _, c := range m.cookies {
		req.AddCookie(c)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting torrents: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v\n", err)
		}
	}()

	var torrents []Torrent
	if err := json.NewDecoder(resp.Body).Decode(&torrents); err != nil {
		return nil, fmt.Errorf("error decoding torrents JSON: %w", err)
	}

	return torrents, nil
}

// Deletes a torrent by hash
func (m *QbitManager) deleteTorrent(hash, name string) {
	data := url.Values{}
	data.Set("hashes", hash)
	data.Set("deleteFiles", "true")

	req, _ := http.NewRequest("POST", m.cfg.QbittorrentUrl+"/api/v2/torrents/delete", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range m.cookies {
		req.AddCookie(c)
	}

	_, err := m.client.Do(req)
	if err != nil {
		log.Printf("Error deleting torrent %s: %v\n", name, err)
	} else {
		log.Printf("[%s] Deleted torrent: %s\n", time.Now().Format(time.RFC3339), name)
	}
}

// Checks if the torrent state indicates it is seeding
func (m *QbitManager) isCompletedState(state string) bool {
	completedStates := []string{
		"completed",
		"uploading",
		"stalledUP",
		"pausedUP",
		"queuedUP",
	}
	return slices.Contains(completedStates, state)
}

// Checks if the torrent hash is already tracked
func (m *QbitManager) isTracked(hash string) bool {
	_, tracked := m.completionTimes[hash]
	return tracked
}

// Checks if the torrent has been completed for longer than the threshold
func (m *QbitManager) completedPastThreshold(hash string) bool {
	if completionTime, exists := m.completionTimes[hash]; exists {
		return time.Since(completionTime) > time.Duration(m.cfg.DeleteAfterMinutes)*time.Minute
	}
	return false
}

// Run starts the polling loop
func (m *QbitManager) Run() {
	log.Println("Starting qBittorrent auto-delete manager...")
	for {
		torrents, err := m.getTorrents()
		if err != nil {
			log.Println("Error getting torrents:", err)
			log.Printf("Retrying in %d seconds...\n", m.cfg.PollIntervalSeconds)
			time.Sleep(time.Duration(m.cfg.PollIntervalSeconds) * time.Second)
			continue
		}

		now := time.Now()
		for _, t := range torrents {
			if m.isCompletedState(t.State) {
				if !m.isTracked(t.Hash) {
					// If the torrent is seeding and not tracked, track it
					log.Printf("[%s] Tracking new seeding torrent: %s\n", now.Format(time.RFC3339), t.Name)
					m.completionTimes[t.Hash] = now
				} else if m.completedPastThreshold(t.Hash) {
					// If tracked, check if it has been seeding long enough and delete it if so
					log.Printf("[%s] Deleting seeding torrent: %s (Hash: %s)\n", now.Format(time.RFC3339), t.Name, t.Hash)
					m.deleteTorrent(t.Hash, t.Name)
					delete(m.completionTimes, t.Hash)
				}
			} else {
				// If the torrent is not seeding, remove it from tracking
				delete(m.completionTimes, t.Hash)
			}
		}

		log.Printf("[%s] Completed checking torrents, sleeping for %d seconds...\n", now.Format(time.RFC3339), m.cfg.PollIntervalSeconds)
		time.Sleep(time.Duration(m.cfg.PollIntervalSeconds) * time.Second)
	}
}
