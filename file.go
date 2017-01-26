package midi

import "io"

type File struct {
	Header FileHeader
	Tracks []Track
}

func Read(from io.Reader) (*File, error) {
	file := new(File)
	header, err := ReadFileHeader(from)
	if err != nil {
		return nil, err
	}
	file.Header = *header

	for i := uint16(0); i < file.Header.NumTrack; i++ {
		track, err := ReadTrack(from)
		if err != nil {
			return nil, err
		}
		file.Tracks = append(file.Tracks, *track)
	}

	return file, nil
}
