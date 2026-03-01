package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type CharBio struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Bio    string `json:"bio"`
	Traits string `json:"traits"`
}

type BioIndex struct {
	bios map[string]CharBio
	keys []string
}

func NewBioIndex(bios map[string]CharBio) BioIndex {
	keys := make([]string, 0, len(bios))
	for key := range bios {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return BioIndex{
		bios: bios,
		keys: keys,
	}
}

func (b BioIndex) Find(query string) (CharBio, bool) {
	value, ok := b.bios[normalizeKey(query)]
	return value, ok
}

func (b BioIndex) List() string {
	return strings.Join(b.keys, ", ")
}

func getBiosFilePath() string {
	path := strings.TrimSpace(os.Getenv("BIOS_FILE"))
	if path == "" {
		return "bios.json"
	}
	return path
}

func loadBios(path string) (map[string]CharBio, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lendo %s: %w", path, err)
	}

	var raw map[string]CharBio
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse JSON %s: %w", path, err)
	}
	if len(raw) == 0 {
		return nil, errors.New("arquivo de bios vazio")
	}

	bios := make(map[string]CharBio, len(raw))
	for key, bio := range raw {
		normalizedKey := normalizeKey(key)
		if normalizedKey == "" {
			return nil, errors.New("encontrei chave vazia em bios.json")
		}
		if bio.Name == "" {
			return nil, fmt.Errorf("bio sem nome para a chave %q", key)
		}
		if _, exists := bios[normalizedKey]; exists {
			return nil, fmt.Errorf("chave duplicada após normalização: %q", normalizedKey)
		}
		bios[normalizedKey] = bio
	}

	return bios, nil
}

func normalizeKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
