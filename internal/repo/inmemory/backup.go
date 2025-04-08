package inmemory

import (
	"bufio"
	"encoding/json"
)

type backupRecord struct {
	OriginalURL string `json:"original_url"`
	ShortID     string `json:"short_id"`
}

func (r *inmemory) restoreBackup() error {
	scanner := bufio.NewScanner(r.file)
	var record backupRecord

	for scanner.Scan() {
		payload := scanner.Bytes()
		err := scanner.Err()
		if err != nil {
			return err
		}

		err = json.Unmarshal(payload, &record)
		if err != nil {
			return err
		}

		r.originalURLs[record.ShortID] = record.OriginalURL
		r.shortIDs[record.OriginalURL] = record.ShortID
	}

	return nil
}

func (r *inmemory) backup(originalURL, shortID string) error {
	record := backupRecord{
		OriginalURL: originalURL,
		ShortID:     shortID,
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}

	payload = append(payload, 10)
	_, err = r.file.Write(payload)
	if err != nil {
		return err
	}

	return nil
}
