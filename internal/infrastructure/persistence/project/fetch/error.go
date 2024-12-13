package fetch

import "errors"

var (
	// ErrCopyingAFileFromLocalToDirWorkingDir represents an error copying a file in the working directory
	ErrCopyingAFileFromLocalToDirWorkingDir = errors.New("An error occurred copying a file in the working directory")
	// ErrCopyingFilesToWorkingDir represents an error copying files to working directory
	ErrCopyingFilesToWorkingDir = errors.New("error copying files to working directory")
	// ErrCreatingADirectoryInLocalDirWorkingDir represents an error creating a directory in the working directory
	ErrCreatingADirectoryInLocalDirWorkingDir = errors.New("An error occurred creating a directory in the working directory")
	// ErrCreatingAFileFromLocalToDirWorkingDir represents an error creating a file in the working directory
	ErrCreatingAFileFromLocalToDirWorkingDir = errors.New("An error occurred creating a file in the working directory")
	// ErrFetchingProjectFromLocalStorage represents an error when fetching a project from local storage
	ErrFetchingProjectFromLocalStorage = errors.New("error fetching a project from local storage")
	// ErrFileSystemNotInitialized represents an error when the filesystem is not initialized
	ErrFileSystemNotInitialized = errors.New("filesystem not initialized")
	// ErrGettingSourceCodeRelativePathFromLocalDir represents an error getting the relative path of the source code
	ErrGettingSourceCodeRelativePathFromLocalDir = errors.New("An error occurred getting the relative path of the source code")
	// ErrInvalidProjectReference represents an error when the project reference is invalid
	ErrInvalidProjectReference = errors.New("invalid project reference")
	// ErrOpeningASourceCodeFileFromLocalDir represents an error opening a source code file
	ErrOpeningASourceCodeFileFromLocalDir = errors.New("An error occurred opening a source code file")
	// ErrProjectNotProvided represents an error when the project is not provided
	ErrProjectNotProvided = errors.New("project not provided")
	// ErrProjectReferenceNotProvided represents an error when the project reference is not provided
	ErrProjectReferenceNotProvided = errors.New("project reference not provided")
	// ErrSourceCodeNotExists represents an error when the source code does not exists
	ErrSourceCodeNotExists = errors.New("source code does not exists")
	// ErrWalkingDirToFetchSourceCodeFromLocalDir represents an error walking through the source code directory
	ErrWalkingDirToFetchSourceCodeFromLocalDir = errors.New("An error occurred walking through the source code directory")
	// ErrWorkingDirNotExists represents an error when the destination to fetch does not exists
	ErrWorkingDirNotExists = errors.New("working directory does not exists")
	// ErrWorkingDirNotProvided represents an error when the working directory is not provided
	ErrWorkingDirNotProvided = errors.New("working directory not provided")
)
