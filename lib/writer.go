package lib

import "io"

// A wrapper around io.Writer that optimizes
// writing in batches
type BatchWriter struct {
	w   io.Writer
	n   int
	err error
}

func NewBatchWriter(w io.Writer) *BatchWriter {
	return &BatchWriter{
		w:   w,
		n:   0,
		err: nil,
	}
}

// Writes buf contents to writer. If error occurs, all subsequent calls do nothing
func (b *BatchWriter) Write(buf []byte) *BatchWriter {
	if b.err != nil {
		return b
	}

	n, err := b.w.Write(buf)
	b.n += n
	b.err = err
	return b
}

// Writes `s` contents to writer. If error occurs, all subsequent calls do nothing
func (b *BatchWriter) WriteString(s string) *BatchWriter {
	if b.err != nil {
		return b
	}

	n, err := io.WriteString(b.w, s)
	b.n += n
	b.err = err
	return b
}

type Formatter interface {
	Format(s string) string
}

// Writes formatted input to the writer
func (b *BatchWriter) WriteFmt(s string, fmt Formatter) *BatchWriter {
	b.WriteString(fmt.Format(s))
	return b
}

// Returns the first error received while writing
func (b *BatchWriter) Err() error {
	return b.err
}

// Returns number of bytes written it total
func (b *BatchWriter) N() int {
	return b.n
}

// Returns the first error received while writing
// Resets the state of buffer
func (b *BatchWriter) Done() error {
	err := b.err
	b.err = nil
	b.n = 0
	return err
}
