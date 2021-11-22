package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	encoding = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// Appends the data in the byte array: p into the store's buffer.
// Returns the number of bytes written: n, the position where the
// bytes are appended from: pos.
func (s *store) Append(p []byte) (n, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	// Write the length of the record
	if err := binary.Write(s.buf, encoding, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	// Write the record
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	// Read the length of the record
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	// Read the record
	b := make([]byte, encoding.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
