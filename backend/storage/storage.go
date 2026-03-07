package storage

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"

	"github.com/mattheweckstein/risk/backend/models"
)

// Store handles saving/loading game states to/from a JSON file.
type Store struct {
	filepath string
	mu       sync.RWMutex
}

// NewStore creates a new Store that persists to the given filepath.
func NewStore(filepath string) *Store {
	return &Store{
		filepath: filepath,
	}
}

// SaveAll writes all games to the JSON file using atomic temp-file + rename.
func (s *Store) SaveAll(games map[string]*models.GameState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(games, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := s.filepath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.filepath)
}

// LoadAll reads all games from the JSON file. If the file does not exist,
// it returns an empty map (not an error).
func (s *Store) LoadAll() (map[string]*models.GameState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return make(map[string]*models.GameState), nil
		}
		return nil, err
	}

	var games map[string]*models.GameState
	if err := json.Unmarshal(data, &games); err != nil {
		return nil, err
	}

	if games == nil {
		games = make(map[string]*models.GameState)
	}

	return games, nil
}
