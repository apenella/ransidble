package local

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

const (
	// ErrGeneratingDataHash is the error message when generating a hash for the record fails.
	ErrGeneratingDataHash = "error generating hash for data"
	// ErrCreatingRecord is the error message when creating a record fails.
	ErrCreatingRecord = "error creating record"
	// ErrGeneratingRecordReferenceInvalidID = "error generating record reference: invalid ID"
	ErrGeneratingRecordReferenceInvalidID = "error generating record reference: invalid ID"
	// ErrGeneratingRecordReferenceInvalidVersion = "error generating record reference: invalid version"
	ErrGeneratingRecordReferenceInvalidVersion = "error generating record reference: invalid version"
	// ErrRecordDataIsEmpty is the error message when the record data is empty.
	ErrRecordDataIsEmpty = "record data is empty"
	// ErrRcordHashIsEmpty is the error message when the record hash is empty.
	ErrRcordHashIsEmpty = "record hash is empty"
	// ErrRecordFailToVerify is the error message when the record fails to verify.
	ErrRecordFailToVerify = "record failed to verify"

	// DefaultRecordVersion is the default version for a record.
	DefaultRecordVersion = "latest"
)

// Record is a struct that represents a record in the local database.
type Record struct {
	// Hash is the hash of the record data.
	Hash string `json:"hash"`
	// CreatedAt is the time when the record was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the record was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// Data is the data of the record.
	Data json.RawMessage `json:"data"`
}

// CreateRecord creates a new instance of Record.
func CreateRecord(id string, data json.RawMessage) (*Record, error) {

	hash, err := generateDataHash(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreatingRecord, err)
	}

	currentTime := time.Now()

	record := &Record{
		Hash:      hash,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Data:      data,
	}

	return record, nil
}

// Verify verifies the record data integrity. Returns true if the record is valid, false otherwise.
func (r *Record) Verify() (bool, error) {

	if r.Data == nil {
		return false, fmt.Errorf(ErrRecordDataIsEmpty)
	}

	if r.Hash == "" {
		return false, fmt.Errorf(ErrRcordHashIsEmpty)
	}

	hash := sha256.Sum256(r.Data)

	if r.Hash != fmt.Sprintf("%x", hash) {
		return false, nil
	}

	return true, nil
}

// generateHash generates a hash for the record.
func generateDataHash(data json.RawMessage) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrGeneratingDataHash, err)
	}

	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash), nil
}

// generateRecordReference generates an homogeneous reference for the record.
func generateRecordReference(id string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("%s. %w", ErrGeneratingRecordReferenceInvalidID, fmt.Errorf("id is empty"))
	}

	hash := sha256.Sum256([]byte(id))

	return fmt.Sprintf("%x", hash), nil
}
