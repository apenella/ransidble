package unpack

import "io"

// TarExtractorer interface used to extract tar files
type TarExtractorer interface {
	Extract(reader io.Reader, dest string) error
}
