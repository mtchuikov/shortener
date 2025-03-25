package inmemory

import (
	"bufio"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type backupRecord struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func (c *cache) backup(url, id string) error {
	const op = "cache.inmemory.cache.backup"

	for id, url := range c.urls {
		record := backupRecord{
			ID:  id,
			URL: url,
		}

		payload, err := jsoniter.Marshal(record)
		if err != nil {
			return fmt.Errorf(
				"%s - %w: %v",
				op, ErrFailedToMarshalBackupRecord, err,
			)
		}

		_, err = c.file.Write(payload)
		if err != nil {
			return fmt.Errorf(
				"%s - %w: %v",
				op, ErrFailedToWriteBackup, err,
			)
		}

		_, err = c.file.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf(
				"%s - %w: %v",
				op, ErrFailedToWriteBackup, err,
			)
		}
	}

	return nil
}

func (c *cache) restoreBackup() error {
	const op = "cache.inmemory.cache.restore_backup"

	scanner := bufio.NewScanner(c.file)

	var record backupRecord
	for scanner.Scan() {
		payload := scanner.Bytes()

		err := scanner.Err()
		if err != nil {
			return fmt.Errorf(
				"%s - %w: %v",
				op, ErrFailedToRestoreBackup, err,
			)
		}

		err = jsoniter.Unmarshal(payload, &record)
		if err != nil {
			return fmt.Errorf(
				"%s - %w: %v",
				op, ErrFailedToRestoreBackup, err,
			)
		}

		c.urls[record.ID] = record.URL
		c.ids[record.URL] = record.ID
	}
	return nil
}
