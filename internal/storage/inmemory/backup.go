package inmemory

import (
	"bufio"
	"encoding/json"
	"fmt"
)

type backupRecord struct {
	Slug        string `json:"slug"`
	OriginalURL string `json:"original_url"`
}

func (s *shortURLs) restoreBackup() error {
	scanner := bufio.NewScanner(s.backupFile)
	var record backupRecord

	for scanner.Scan() {
		payload := scanner.Bytes()
		err := scanner.Err()
		if err != nil {
			err = fmt.Errorf("failed to restore backup: %w", err)
			return err
		}

		err = json.Unmarshal(payload, &record)
		if err != nil {
			err = fmt.Errorf("failed to restore backup: %w", err)
			return err
		}

		s.slugs[record.OriginalURL] = record.Slug
		s.originalURLs[record.Slug] = record.OriginalURL
	}

	return nil
}

func (s *shortURLs) backup(slug, originalURL string) error {
	payload, err := json.Marshal(
		backupRecord{
			Slug:        slug,
			OriginalURL: originalURL,
		})
	if err != nil {
		err = fmt.Errorf("failed to backup: %w", err)
		return err
	}

	payload = append(payload, 10)
	_, err = s.backupFile.Write(payload)
	if err != nil {
		err = fmt.Errorf("failed to backup: %w", err)
		return err
	}

	return nil
}
