package inmemory

import "errors"

var (
	ErrFailedToOpenBackupFile      = errors.New("failed to open backup file")
	ErrFailedToRestoreBackup       = errors.New("failed to restore backup")
	ErrFailedToWriteBackup         = errors.New("failed to write backup record")
	ErrFailedToMarshalBackupRecord = errors.New("failed to marshal backup record")
)
