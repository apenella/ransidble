package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

const (
	// ErrDatabasePathNotInitialized is the error message when the database path is not initialized.
	ErrDatabasePathNotInitialized = "database path not initialized"
	// ErrDataToWriteIsNotProvided is the error message when the data to write is not provided.
	ErrDataToWriteIsNotProvided = "data to write is not provided"
	// ErrExtractingProjectFromRecord is the error message when extracting the project from the record fails.
	ErrExtractingProjectFromRecord = "error extracting project from record"
	// ErrExtractingRecordFromDatabase is the error message when extracting the record from the database fails.
	ErrExtractingRecordFromDatabase = "error extracting record from database"
	// ErrFilesystemNotInitialized is the error message when the filesystem is not initialized.
	ErrFilesystemNotInitialized = "filesystem not initialized"
	// ErrIDIsNotProvided is the error message when the ID is required.
	ErrIDIsNotProvided = "ID is not provided"
	// ErrInvalidRecordFormat is the error message when the record format is invalid.
	ErrInvalidRecordFormat = "invalid record format"
	// ErrInvalidRecordFormatIsDir is the error message when the record format is a directory.
	ErrInvalidRecordFormatIsDir = "record format is a directory"
	// ErrMarshalingDataToGenerateRecord is the error message when marshaling the data to write fails.
	ErrMarshalingDataToGenerateRecord = "error marshaling data to generate the record"
	// ErrMarshalingRecordToWrite is the error message when marshaling the record to write fails.
	ErrMarshalingRecordToWrite = "error marshaling record to write"
	// ErrOpeningFileToWriteRecord is the error message when opening the file to write the record fails.
	ErrOpeningFileToWriteRecord = "error opening file to write record"
	// ErrReadingRecord is the error message when reading the record fails.
	ErrReadingRecord = "error reading record"
	// ErrReadingRecordNotFound is the error message when the record is not found.
	ErrReadingRecordNotFound = "record not found"
	// ErrReadingRecordsFromDatabase is the error message when reading the records from the database fails.
	ErrReadingRecordsFromDatabase = "error reading records from database"
	// ErrVerifyingRecord is the error message when verifying the record fails.
	ErrVerifyingRecord = "error verifying record"
	// ErrVerifyingRecordInvalidHash is the error message when the record hash is invalid.
	ErrVerifyingRecordInvalidHash = "record hash is invalid"
	// ErrWritingRecord is the error message when writing the record fails.
	ErrWritingRecord = "error writing record"
	// ErrDatabasePathMustExist is the error message when the database path must exist.
	ErrDatabasePathMustExist = "database path must exist"
	// ErrRemovingRecord is the error message when removing the record fails.
	ErrRemovingRecord = "error removing record"
	// ErrProjectExists is the error message when the project already exists.
	ErrProjectExists = "error project already exists"
	// ErrStoringProject is the error message when storing the project fails.
	ErrStoringProject = "error storing project"
)

// DatabaseDriver is a struct that represents a local database to persist the projects references.
type DatabaseDriver struct {
	// fs path where projects are stored
	fs afero.Fs
	// path is the path to the directory where the projects references are stored
	path string

	logger repository.Logger
}

// NewDatabaseDriver creates a new instance of DatabaseDriver.
func NewDatabaseDriver(fs afero.Fs, path string, logger repository.Logger) *DatabaseDriver {
	return &DatabaseDriver{
		fs:     fs,
		logger: logger,
		path:   path,
	}
}

// Find a project from the local database.
func (db *DatabaseDriver) Find(id string) (*entity.Project, error) {
	return db.read(id)
}

// FindAll reads all projects from the local database.
func (db *DatabaseDriver) FindAll() ([]*entity.Project, error) {
	return db.readAll()
}

// Store stores a project in the local database.
func (db *DatabaseDriver) Store(id string, data *entity.Project) error {
	return db.write(id, data)
}

// SafeStore stores a project in the local database.
func (db *DatabaseDriver) SafeStore(id string, data *entity.Project) error {

	exists, err := db.exists(id)

	if err != nil {
		return fmt.Errorf("%s: %s", ErrStoringProject, err.Error())
	}

	if exists {
		return fmt.Errorf("%s: %s %s", ErrStoringProject, id, ErrProjectExists)
	}

	return db.write(id, data)
}

// Delete deletes a project from the local database.
func (db *DatabaseDriver) Delete(id string) error {
	return db.remove(id)
}

// Private methods to manage the local database

// Read a project from the local database.
func (db *DatabaseDriver) read(id string) (*entity.Project, error) {

	var err error
	var recordContent []byte
	var recordFile afero.File
	var recordFileInfo os.FileInfo
	var project *entity.Project
	var record *Record
	var verifiedRecord bool

	if id == "" {
		db.logger.Error(
			ErrIDIsNotProvided,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
			},
		)

		return nil, fmt.Errorf("%s", ErrIDIsNotProvided)
	}

	if db.fs == nil {
		db.logger.Error(
			ErrFilesystemNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf("%s", ErrFilesystemNotInitialized)
	}

	if db.path == "" {
		db.logger.Error(
			ErrDatabasePathNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf("%s", ErrDatabasePathNotInitialized)
	}

	recordPath := filepath.Join(db.path, id)

	recordFileInfo, err = db.fs.Stat(recordPath)
	if err != nil {
		if os.IsNotExist(err) {
			msgErr := fmt.Sprintf("%s: %s: %s", ErrReadingRecord, ErrReadingRecordNotFound, err.Error())
			db.logger.Error(
				msgErr,
				map[string]interface{}{
					"component": "DatabaseDriver.Read",
					"package":   packageName,
					"record_id": id,
				},
			)

			return nil, fmt.Errorf(msgErr)
		}

		msgErr := fmt.Sprintf("%s: %s", ErrReadingRecord, err.Error())
		db.logger.Error(
			msgErr,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf(msgErr)
	}

	if recordFileInfo.IsDir() {
		msgErr := fmt.Sprintf("%s: %s", ErrInvalidRecordFormat, ErrInvalidRecordFormatIsDir)
		db.logger.Error(
			msgErr,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf(msgErr)
	}

	recordFile, err = db.fs.Open(recordPath)
	if err != nil {
		msgErr := fmt.Sprintf("%s: %s", ErrReadingRecord, err.Error())
		db.logger.Error(
			msgErr,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf(msgErr)
	}
	defer recordFile.Close()

	recordContent, err = io.ReadAll(recordFile)
	if err != nil {
		msgErr := fmt.Sprintf("%s: %s", ErrReadingRecord, err.Error())
		db.logger.Error(
			msgErr,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf(msgErr)
	}

	err = json.Unmarshal(recordContent, &record)
	if err != nil {
		msgErr := fmt.Sprintf("%s: %s", ErrExtractingRecordFromDatabase, err.Error())
		db.logger.Error(
			msgErr,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf(msgErr)
	}

	verifiedRecord, err = record.Verify()
	if err != nil {
		msgErr := fmt.Sprintf("%s: %s", ErrVerifyingRecord, err.Error())
		db.logger.Error(
			msgErr,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf(msgErr)
	}

	if !verifiedRecord {
		msgErr := fmt.Sprintf("%s: %s", ErrVerifyingRecord, ErrVerifyingRecordInvalidHash)
		db.logger.Error(
			msgErr,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf(msgErr)
	}

	err = json.Unmarshal(record.Data, &project)
	if err != nil {
		msgErr := fmt.Sprintf("%s: %s", ErrExtractingProjectFromRecord, err.Error())
		db.logger.Error(
			msgErr,
			map[string]interface{}{
				"component": "DatabaseDriver.Read",
				"package":   packageName,
				"record_id": id,
			},
		)

		return nil, fmt.Errorf(msgErr)
	}

	return project, nil
}

// ReadAll reads all projects from the local database.
func (db *DatabaseDriver) readAll() ([]*entity.Project, error) {

	var projectList []*entity.Project
	var err error

	if db.fs == nil {
		db.logger.Error(
			ErrFilesystemNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.ReadAll",
				"package":   packageName,
			},
		)

		return nil, fmt.Errorf("%s", ErrFilesystemNotInitialized)
	}
	if db.path == "" {
		db.logger.Error(
			ErrDatabasePathNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.ReadAll",
				"package":   packageName,
			},
		)

		return nil, fmt.Errorf("%s", ErrDatabasePathNotInitialized)
	}

	err = afero.Walk(db.fs, db.path, func(path string, info os.FileInfo, err error) error {
		var project *entity.Project
		var errIteration error

		if err != nil {
			db.logger.Error(
				fmt.Sprintf("%s: %s", ErrReadingRecordsFromDatabase, err.Error()),
				map[string]interface{}{
					"component": "DatabaseDriver.ReadAll",
					"package":   packageName,
					"record_id": info.Name(),
				},
			)
			return fmt.Errorf("%s: %w", ErrReadingRecordsFromDatabase, err)
		}

		if info.IsDir() {
			db.logger.Debug(
				fmt.Sprintf("%s: %s", ErrInvalidRecordFormat, ErrInvalidRecordFormatIsDir),
				map[string]interface{}{
					"component": "DatabaseDriver.ReadAll",
					"package":   packageName,
					"record_id": info.Name(),
				},
			)
			return nil
		}

		project, errIteration = db.read(info.Name())
		if errIteration != nil {
			db.logger.Error(
				fmt.Sprintf("%s: %s", ErrReadingRecordsFromDatabase, errIteration.Error()),
				map[string]interface{}{
					"component": "DatabaseDriver.ReadAll",
					"package":   packageName,
					"record_id": info.Name(),
				},
			)

			return nil
		}

		projectList = append(projectList, project)

		return nil
	})

	if err != nil {
		db.logger.Error(
			ErrReadingRecordsFromDatabase,
			map[string]interface{}{
				"component": "DatabaseDriver.ReadAll",
				"package":   packageName,
				"record_id": err.Error(),
			},
		)
		return nil, fmt.Errorf("%s: %w", ErrReadingRecordsFromDatabase, err)
	}

	return projectList, nil
}

// Write a project to the local database.
func (db *DatabaseDriver) write(id string, data *entity.Project) error {

	var err error
	var record *Record
	var recordFile afero.File

	if id == "" {
		db.logger.Error(
			ErrIDIsNotProvided,
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
			},
		)

		return fmt.Errorf("%s", ErrIDIsNotProvided)
	}

	if data == nil {
		db.logger.Error(
			ErrDataToWriteIsNotProvided,
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"record_id": id,
				"package":   packageName,
			},
		)
		return fmt.Errorf("%s", ErrDataToWriteIsNotProvided)
	}

	if db.fs == nil {
		db.logger.Error(
			ErrFilesystemNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
				"project":   data,
				"record_id": id,
			},
		)
		return fmt.Errorf("%s", ErrFilesystemNotInitialized)
	}

	if db.path == "" {
		db.logger.Error(
			ErrDatabasePathNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
				"project":   data,
				"record_id": id,
			},
		)
		return fmt.Errorf("%s", ErrDatabasePathNotInitialized)
	}

	_, err = db.fs.Stat(db.path)
	if err != nil {
		db.logger.Error(
			fmt.Sprintf("%s: %s", ErrDatabasePathMustExist, err.Error()),
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
				"path":      db.path,
				"project":   data,
				"record_id": id,
			},
		)
		return fmt.Errorf("%s: %w", ErrDatabasePathMustExist, err)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		db.logger.Error(
			fmt.Sprintf("%s: %s", ErrMarshalingDataToGenerateRecord, err.Error()),
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
				"record_id": id,
				"project":   data,
			},
		)
		return fmt.Errorf("%s: %w", ErrMarshalingDataToGenerateRecord, err)
	}

	record, err = CreateRecord(id, dataBytes)
	if err != nil {
		db.logger.Error(
			fmt.Sprintf("%s: %s", ErrCreatingRecord, err.Error()),
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
				"record_id": id,
				"project":   data,
			},
		)
	}

	recordFile, err = db.fs.OpenFile(
		filepath.Join(db.path, id),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0644,
	)

	if err != nil {
		db.logger.Error(
			fmt.Sprintf("%s: %s", ErrOpeningFileToWriteRecord, err.Error()),
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
				"record_id": id,
				"project":   data,
			},
		)
		return fmt.Errorf("%s: %w", ErrOpeningFileToWriteRecord, err)
	}
	defer recordFile.Close()

	recordContent, err := json.Marshal(record)
	if err != nil {
		db.logger.Error(
			fmt.Sprintf("%s: %s", ErrMarshalingRecordToWrite, err.Error()),
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
				"record_id": id,
				"project":   data,
			},
		)
		return fmt.Errorf("%s: %w", ErrMarshalingRecordToWrite, err)
	}
	_, err = recordFile.Write(recordContent)
	if err != nil {
		db.logger.Error(
			fmt.Sprintf("%s: %s", ErrWritingRecord, err.Error()),
			map[string]interface{}{
				"component": "DatabaseDriver.Write",
				"package":   packageName,
				"record_id": id,
				"project":   data,
			},
		)
		return fmt.Errorf("%s: %w", ErrWritingRecord, err)
	}

	return nil
}

// Remove a project from the local database.
func (db *DatabaseDriver) remove(id string) error {
	var reference string
	var err error
	var exists bool

	if id == "" {
		db.logger.Error(
			ErrIDIsNotProvided,
			map[string]interface{}{
				"component": "DatabaseDriver.Remove",
				"package":   packageName,
			},
		)

		return fmt.Errorf("%s", ErrIDIsNotProvided)
	}

	if db.fs == nil {
		db.logger.Error(
			ErrFilesystemNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.Remove",
				"package":   packageName,
				"record_id": id,
			},
		)

		return fmt.Errorf("%s", ErrFilesystemNotInitialized)
	}

	if db.path == "" {
		db.logger.Error(
			ErrDatabasePathNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.Remove",
				"package":   packageName,
				"record_id": id,
			},
		)
		return fmt.Errorf("%s", ErrDatabasePathNotInitialized)

	}

	exists, err = db.exists(id)
	if !exists {
		db.logger.Debug(
			"Record could not be removed: not found",
			map[string]interface{}{
				"component": "DatabaseDriver.Remove",
				"package":   packageName,
				"record_id": id,
			},
		)
		return nil
	}
	if err != nil {
		db.logger.Error(
			fmt.Sprintf("%s: %s", ErrRemovingRecord, err.Error()),
			map[string]interface{}{
				"component": "DatabaseDriver.Remove",
				"package":   packageName,
				"record_id": id,
			},
		)

		return fmt.Errorf("%s: %w", ErrRemovingRecord, err)
	}

	reference = filepath.Join(db.path, id)
	err = db.fs.Remove(reference)
	if err != nil {
		db.logger.Error(
			fmt.Sprintf("%s: %s", ErrRemovingRecord, err.Error()),
			map[string]interface{}{
				"component": "DatabaseDriver.Remove",
				"package":   packageName,
				"record_id": id,
			},
		)

		return fmt.Errorf("%s: %w", ErrRemovingRecord, err)
	}

	db.logger.Debug(
		"Record removed",
		map[string]interface{}{
			"component": "DatabaseDriver.Remove",
			"package":   packageName,
			"record_id": id,
		},
	)

	return nil
}

// Exists checks if a project exists in the local database.
func (db *DatabaseDriver) exists(id string) (bool, error) {
	var reference string
	var err error

	if id == "" {
		db.logger.Error(
			ErrIDIsNotProvided,
			map[string]interface{}{
				"component": "DatabaseDriver.Exists",
				"package":   packageName,
				"record_id": id,
			},
		)

		return false, fmt.Errorf("%s", ErrIDIsNotProvided)
	}

	if db.fs == nil {
		db.logger.Error(
			ErrFilesystemNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.Exists",
				"package":   packageName,
				"record_id": id,
			},
		)
		return false, fmt.Errorf("%s", ErrFilesystemNotInitialized)
	}

	if db.path == "" {
		db.logger.Error(
			ErrDatabasePathNotInitialized,
			map[string]interface{}{
				"component": "DatabaseDriver.Exists",
				"package":   packageName,
				"record_id": id,
			},
		)
		return false, fmt.Errorf("%s", ErrDatabasePathNotInitialized)
	}

	reference = filepath.Join(db.path, id)
	_, err = db.fs.Stat(reference)
	if err != nil {
		if os.IsNotExist(err) {
			db.logger.Debug(
				ErrReadingRecordNotFound,
				map[string]interface{}{
					"component": "DatabaseDriver.Exists",
					"package":   packageName,
					"record_id": id,
				},
			)

			return false, nil
		}

		db.logger.Error(
			ErrReadingRecord,
			map[string]interface{}{
				"component": "DatabaseDriver.Exists",
				"package":   packageName,
				"record_id": id,
			},
		)

		return false, fmt.Errorf("%s: %w", ErrReadingRecord, err)
	}

	return true, nil
}
