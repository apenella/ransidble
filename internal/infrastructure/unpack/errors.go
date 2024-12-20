package unpack

import "errors"

var (
	// ErrCreatingWorkingDirFolder represents an error creating working directory folder
	ErrCreatingWorkingDirFolder = errors.New("error creating working directory folder")
	// ErrCopyingFilesInWorkingDir represents an error copying files in working directory
	ErrCopyingFilesInWorkingDir = errors.New("error copying files in working directory")
	// ErrRemovingWorkingDirFolder represents an error removing working directory folder
	ErrRemovingWorkingDirFolder = errors.New("error removing working directory folder")
	// ErrProjectNotProvided represents an error when the project is not provided
	ErrProjectNotProvided = errors.New("project not provided")
	// ErrWorkingDirNotProvided represents an error when the working directory is not provided
	ErrWorkingDirNotProvided = errors.New("working directory not provided")
	// ErrWorkingDirIsNotDirectory represents an error when the working directory is not a directory
	ErrWorkingDirIsNotDirectory = errors.New("working directory is not a directory")
	// ErrDescribingWorkingDir represents an error when the working directory cannot be described
	ErrDescribingWorkingDir = errors.New("an error occurred describing working directory")
	// ErrFilesystemNotProvided represents an error when the filesystem is not provided
	ErrFilesystemNotProvided = errors.New("filesystem not provided")
	// ErrWorkingDirNotExists represents an error when the destination to fetch does not exists
	ErrWorkingDirNotExists = errors.New("working directory does not exists")
	// ErrTarExtractorNotProvided is returned when the tar extractor is not provided
	ErrTarExtractorNotProvided = errors.New("Tar extractor not provided")
	// ErrSourceCodeFileNotExist is returned when the source code file does not exist
	ErrSourceCodeFileNotExist = errors.New("source code file does not exist")
	// ErrOpeningSourceCodeFile is returned when the source code file cannot be opened
	ErrOpeningSourceCodeFile = errors.New("an error occurred opening source code file")
	// ErrCreatingGzipReader is returned when the gzip reader cannot be created
	ErrCreatingGzipReader = errors.New("an error occurred creating gzip reader")
	// ErrExtractingSourceCodeFile is returned when the source code file cannot be extracted
	ErrExtractingSourceCodeFile = errors.New("an error occurred extracting source code file")
	// ErrDescribingProjectReferenece is returned when the project reference cannot be described
	ErrDescribingProjectReferenece = errors.New("an error occurred describing project reference")
	// ErrProjectReferenceNotProvided is returned when the project reference is not provided
	ErrProjectReferenceNotProvided = errors.New("project reference not provided")
)
