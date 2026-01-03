package jsonld

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// StreamError represents an error that occurred during streaming serialization.
type StreamError struct {
	Index   int
	Message string
	Err     error
}

func (e *StreamError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("stream error at index %d: %s: %v", e.Index, e.Message, e.Err)
	}
	return fmt.Sprintf("stream error at index %d: %s", e.Index, e.Message)
}

func (e *StreamError) Unwrap() error {
	return e.Err
}

// StreamResult contains the result of a streaming operation.
type StreamResult struct {
	RecordsWritten int
	Errors         []*StreamError
}

// HasErrors returns true if any errors occurred during streaming.
func (r *StreamResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// Buffer pool configuration constants.
const (
	// defaultBufferSize is the initial capacity for pooled buffers (4KB).
	defaultBufferSize = 4096
	// maxPooledBufferSize is the maximum buffer size to return to pool (64KB).
	maxPooledBufferSize = 65536
)

// bufferPool is a sync.Pool for reusing byte buffers during streaming.
//
//nolint:gochecknoglobals // Intentional optimization for buffer reuse
var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, defaultBufferSize))
	},
}

func getBuffer() *bytes.Buffer {
	buf, ok := bufferPool.Get().(*bytes.Buffer)
	if !ok {
		return bytes.NewBuffer(make([]byte, 0, defaultBufferSize))
	}
	buf.Reset()
	return buf
}

func putBuffer(buf *bytes.Buffer) {
	if buf.Cap() > maxPooledBufferSize {
		// Don't pool very large buffers
		return
	}
	buf.Reset() // Reset before pooling (defensive programming)
	bufferPool.Put(buf)
}

// SerializeStream writes multiple FocusCostRecords to a writer as a JSON-LD array.
//
// The output format is a JSON array: [record1, record2, ...]
// Memory usage is bounded - only one record is held in memory at a time.
//
// Error Handling:
//   - Returns immediately if the opening bracket cannot be written
//   - Serialization errors (invalid records) are collected in StreamResult.Errors
//     and processing continues with the next record
//   - Write errors are also collected in StreamResult.Errors, allowing partial
//     output when the writer recovers (e.g., temporary network issues)
//   - Returns immediately if the closing bracket cannot be written
//   - Call StreamResult.HasErrors() to check if any errors occurred
func (s *Serializer) SerializeStream(records <-chan *pbc.FocusCostRecord, w io.Writer) (*StreamResult, error) {
	result := &StreamResult{
		Errors: make([]*StreamError, 0),
	}

	// Write opening bracket
	if _, err := w.Write([]byte("[\n")); err != nil {
		return result, fmt.Errorf("failed to write opening bracket: %w", err)
	}

	index := 0
	first := true

	for record := range records {
		// Get a buffer from the pool
		buf := getBuffer()

		// Serialize the record
		data, err := s.Serialize(record)
		if err != nil {
			result.Errors = append(result.Errors, &StreamError{
				Index:   index,
				Message: "serialization failed",
				Err:     err,
			})
			putBuffer(buf)
			index++
			continue
		}

		// Write separator if not first record
		if !first {
			buf.WriteString(",\n")
		}
		first = false

		// Indent the record for readability in streaming output
		if s.options.PrettyPrint {
			var indented bytes.Buffer
			if indentErr := json.Indent(&indented, data, "  ", "  "); indentErr == nil {
				buf.WriteString("  ")
				buf.Write(indented.Bytes())
			} else {
				buf.Write(data)
			}
		} else {
			buf.Write(data)
		}

		// Write to output
		if _, writeErr := w.Write(buf.Bytes()); writeErr != nil {
			result.Errors = append(result.Errors, &StreamError{
				Index:   index,
				Message: "write failed",
				Err:     writeErr,
			})
			putBuffer(buf)
			index++
			continue
		}

		putBuffer(buf)
		result.RecordsWritten++
		index++
	}

	// Write closing bracket
	if _, err := w.Write([]byte("\n]")); err != nil {
		return result, fmt.Errorf("failed to write closing bracket: %w", err)
	}

	return result, nil
}

// SerializeSlice serializes a slice of FocusCostRecords to JSON-LD array.
//
// This is a convenience method that creates a channel from the slice
// and calls SerializeStream.
func (s *Serializer) SerializeSlice(records []*pbc.FocusCostRecord, w io.Writer) (*StreamResult, error) {
	ch := make(chan *pbc.FocusCostRecord, len(records))
	for _, r := range records {
		ch <- r
	}
	close(ch)
	return s.SerializeStream(ch, w)
}

// SerializeBatch serializes multiple records and returns the complete JSON-LD array.
//
// Unlike SerializeStream, this method buffers all output in memory.
// Use SerializeStream for very large datasets.
func (s *Serializer) SerializeBatch(records []*pbc.FocusCostRecord) ([]byte, *StreamResult, error) {
	var buf bytes.Buffer
	result, err := s.SerializeSlice(records, &buf)
	if err != nil {
		return nil, result, err
	}
	return buf.Bytes(), result, nil
}
