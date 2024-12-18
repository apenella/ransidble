package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

var (
	// ErrUnableToUntar is returned when the tar file is unable to be untarred
	ErrUnableToUntar = fmt.Errorf("unable to untar")
	// ErrExtractingFileFromTar is returned when the file is unable to be extracted from the tar file
	ErrExtractingFileFromTar = fmt.Errorf("unable to extract file from tar")
	// ErrTarReading is returned when the tar file is unable to be read
	ErrTarReading = fmt.Errorf("unable to read tar file")
	// ErrCopyingContentFromTar is returned when the content is unable to be copied from the tar file
	ErrCopyingContentFromTar = fmt.Errorf("unable to copy content from tar")
	// ErrCreatingFileFromTar is returned when the file is unable to be created from the tar file
	ErrCreatingFileFromTar = fmt.Errorf("unable to create file from tar")
	// ErrReaderNotProvided is returned when the reader is not provided
	ErrReaderNotProvided = fmt.Errorf("reader not provided")
	// ErrDestinationNotProvided is returned when the destination is not provided
	ErrDestinationNotProvided = fmt.Errorf("destination not provided")
	// ErrFilesystemNotProvided is returned when the filesystem is not provided
	ErrFilesystemNotProvided = fmt.Errorf("filesystem not provided")
	// ErrTarFileHeaderNotProvided is returned when the tar file header is not provided
	ErrTarFileHeaderNotProvided = fmt.Errorf("tar file header not provided")
)

// Tar is a struct that implements the Tar operations
type Tar struct {
	// fs is the filesystem
	fs afero.Fs
	// logger is a logger interface
	logger repository.Logger
}

// NewTar creates a new Tar struct
func NewTar(fs afero.Fs, logger repository.Logger) *Tar {
	return &Tar{
		fs:     fs,
		logger: logger,
	}
}

// Extract untar the io.Reader into the destination
func (t *Tar) Extract(r io.Reader, destination string) error {

	if r == nil {
		t.logger.Error(
			ErrReaderNotProvided.Error(),
			map[string]interface{}{
				"component": "Tar.Extract",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
			})
		return ErrReaderNotProvided
	}

	if destination == "" {
		t.logger.Error(
			ErrDestinationNotProvided.Error(),
			map[string]interface{}{
				"component": "Tar.Extract",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
			})
		return ErrDestinationNotProvided
	}

	if t.fs == nil {
		t.logger.Error(
			ErrFilesystemNotProvided.Error(),
			map[string]interface{}{
				"component": "Tar.Extract",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
			})
		return ErrFilesystemNotProvided
	}

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}

		if err != nil {
			t.logger.Error(
				fmt.Sprintf("%s: %s", ErrTarReading, err),
				map[string]interface{}{
					"component": "Tar.Extract",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
				})
			return fmt.Errorf("%s: %w", ErrTarReading, err)
		}

		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err := t.fs.MkdirAll(target, os.FileMode(header.Mode))
			if err != nil {
				t.logger.Error(
					fmt.Sprintf("%s: %s", ErrUnableToUntar, err),
					map[string]interface{}{
						"component": "Tar.Extract",
						"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
						"file":      header.Name,
						"type":      "directory",
					})
				return fmt.Errorf("%s: %w", ErrUnableToUntar, err)
			}
		case tar.TypeReg:
			err := t.extractRegularFile(tr, header, target)
			if err != nil {
				t.logger.Error(
					fmt.Sprintf("%s: %s", ErrExtractingFileFromTar, err),
					map[string]interface{}{
						"component": "Tar.Extract",
						"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
						"file":      header.Name,
						"type":      "file",
					})
				return fmt.Errorf("%s: %w", ErrExtractingFileFromTar, err)
			}
		// The following cases are not supported yet, it would be required when fetching files from a git source
		// https://github.com/golang/build/blob/master/internal/untar/untar.go#L131C1-L132C48
		// case tar.TypeXGlobalHeader:
		// 	// git archive generates these. Ignore them.
		default:
			t.logger.Error(
				ErrUnableToUntar.Error(),
				map[string]interface{}{
					"component": "Tar.Extract",
					"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
					"file":      header.Name,
					"type":      string(header.Typeflag),
				})

			return ErrUnableToUntar
		}
	}

	return nil
}

// extractRegularFile extracts a regular file from the tar file
func (t *Tar) extractRegularFile(tr *tar.Reader, header *tar.Header, destination string) (err error) {
	var file afero.File

	if header == nil {
		t.logger.Error(
			ErrTarFileHeaderNotProvided.Error(),
			map[string]interface{}{
				"component": "Tar.extractRegularFile",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
			})
		return ErrTarFileHeaderNotProvided
	}

	if destination == "" {
		t.logger.Error(
			ErrDestinationNotProvided.Error(),
			map[string]interface{}{
				"component": "Tar.extractRegularFile",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
				"file":      header.Name,
				"type":      "file",
			})
		return ErrDestinationNotProvided
	}

	if tr == nil {
		t.logger.Error(
			ErrReaderNotProvided.Error(),
			map[string]interface{}{
				"component":   "Tar.extractRegularFile",
				"destination": destination,
				"file":        header.Name,
				"package":     "github.com/apenella/ransidble/internal/infrastructure/tar",
				"type":        "file",
			})
		return ErrReaderNotProvided
	}

	file, err = t.fs.OpenFile(destination, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
	defer func() {
		err = file.Close()
	}()

	if err != nil {
		t.logger.Error(
			fmt.Sprintf("%s: %s", ErrCreatingFileFromTar, err),
			map[string]interface{}{
				"component": "Tar.extractRegularFile",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
				"file":      header.Name,
				"type":      "file",
			})

		return fmt.Errorf("%s: %w", ErrCreatingFileFromTar, err)
	}

	if _, err = io.Copy(file, tr); err != nil {
		t.logger.Error(
			fmt.Sprintf("%s: %s", ErrCopyingContentFromTar, err),
			map[string]interface{}{
				"component": "Tar.extractRegularFile",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/tar",
				"file":      header.Name,
				"type":      "file",
			})
		return fmt.Errorf("%s: %w", ErrCopyingContentFromTar, err)
	}

	return nil
}
