package inmemory

import (
	"bufio"
	"encoding/json"
	"fmt"
)

type backupRecord struct {
	OriginalURL string `json:"original_url"`
	ShortID     string `json:"short_id"`
}

const wFailedToBackupRecord = "failed to backup record: %w"

func (r *repo) backup(originalURL, shortID string) error {
	record := backupRecord{
		OriginalURL: originalURL,
		ShortID:     shortID,
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf(wFailedToBackupRecord, err)
	}

	payload = append(payload, 10)
	_, err = r.file.Write(payload)
	if err != nil {
		return fmt.Errorf(wFailedToBackupRecord, err)
	}

	return nil
}

const wFailedToRestoreRecord = "failed to restore record: %w"

func (r *repo) restoreBackup() error {
	scanner := bufio.NewScanner(r.file)
	var record backupRecord

	for scanner.Scan() {
		payload := scanner.Bytes()
		err := scanner.Err()
		if err != nil {
			return fmt.Errorf(wFailedToRestoreRecord, err)
		}

		err = json.Unmarshal(payload, &record)
		if err != nil {
			return fmt.Errorf(wFailedToRestoreRecord, err)
		}

		r.originalURLs[record.ShortID] = record.OriginalURL
		r.shortIDs[record.OriginalURL] = record.ShortID
	}

	return nil
}
