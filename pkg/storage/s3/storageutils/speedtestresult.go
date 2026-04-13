package storageutils

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/bytesutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type SpeedTestResult struct {
	Size       int64
	ChunkSize  int64
	FileName   string
	WriteSpeed float64
	ReadSpeed  float64
}

func (s *SpeedTestResult) GetSize() (int64, error) {
	if s.Size < 0 {
		return 0, tracederrors.TracedErrorf("Invalid size: %d", s.Size)
	}

	return s.Size, nil
}

func (s *SpeedTestResult) GetSizeHumanReadable() (string, error) {
	size, err := s.GetSize()
	if err != nil {
		return "", err
	}

	return bytesutils.GetSizeAsHumanReadableString(size)
}

func (s *SpeedTestResult) GetChunkSize() (int64, error) {
	if s.ChunkSize < 0 {
		return 0, tracederrors.TracedErrorf("Invalid chunk size: %d", s.ChunkSize)
	}

	return s.ChunkSize, nil
}

func (s *SpeedTestResult) GetChunkSizeHumanReadable() (string, error) {
	chunkSize, err := s.GetChunkSize()
	if err != nil {
		return "", err
	}

	return bytesutils.GetSizeAsHumanReadableString(chunkSize)
}

func (s *SpeedTestResult) GetFileName() (string, error) {
	if s.FileName == "" {
		return "", tracederrors.TracedErrorf("FileName is not set")
	}

	return s.FileName, nil
}

func (s *SpeedTestResult) GetWriteSpeed() (float64, error) {
	if s.WriteSpeed < 0 {
		return 0, tracederrors.TracedErrorf("Invalid write speed: %f", s.WriteSpeed)
	}

	return s.WriteSpeed, nil
}

func (s *SpeedTestResult) GetWriteSpeedHumanReadable() (string, error) {
	speed, err := s.GetWriteSpeed()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2fMB/s", speed/1024/1024), nil
}

func (s *SpeedTestResult) GetReadSpeed() (float64, error) {
	if s.ReadSpeed < 0 {
		return 0, tracederrors.TracedErrorf("Invalid read speed: %f", s.ReadSpeed)
	}

	return s.ReadSpeed, nil
}

func (s *SpeedTestResult) GetReadSpeedHumanReadable() (string, error) {
	speed, err := s.GetReadSpeed()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2fMB/s", speed/1024/1024), nil
}

func (s *SpeedTestResult) GetResultMessage() (string, error) {
	size, err := s.GetSizeHumanReadable()
	if err != nil {
		return "", err
	}

	chunkSize, err := s.GetChunkSizeHumanReadable()
	if err != nil {
		return "", err
	}

	fileName, err := s.GetFileName()
	if err != nil {
		return "", err
	}

	writeSpeed, err := s.GetWriteSpeedHumanReadable()
	if err != nil {
		return "", err
	}

	readSpeed, err := s.GetReadSpeedHumanReadable()
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf("Speedtest for file '%s' with size '%s' (chunkSize='%s'): WriteSpeed = '%s', ReadSpeed = '%s'.", fileName, size, chunkSize, writeSpeed, readSpeed)

	return message, nil
}
