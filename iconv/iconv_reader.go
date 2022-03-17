package iconv

import (
	"io"
	"syscall"
)

type Reader struct {
	source            io.Reader
	converter         *Converter
	buffer            []byte
	readPos, writePos int
	err               error
}

func NewReader(source io.Reader, fromEncoding string, toEncoding string) (*Reader, error) {
	// create a converter
	converter, err := NewConverter(fromEncoding, toEncoding)
	if err == nil {
		return NewReaderFromConverter(source, converter), err
	}
	// return the error
	return nil, err
}

func NewReaderFromConverter(source io.Reader, converter *Converter) (reader *Reader) {
	reader = new(Reader)
	// copy elements
	reader.source = source
	reader.converter = converter
	// create 8K buffers
	reader.buffer = make([]byte, 8*1024)
	return reader
}

func (c *Reader) fillBuffer() {
	// slide existing data to beginning
	if c.readPos > 0 {
		// copy current bytes - is this guaranteed safe?
		copy(c.buffer, c.buffer[c.readPos:c.writePos])
		// adjust positions
		c.writePos -= c.readPos
		c.readPos = 0
	}
	// read new data into buffer at write position
	bytesRead, err := c.source.Read(c.buffer[c.writePos:])
	// adjust write position
	c.writePos += bytesRead
	// track any reader error / EOF
	if err != nil {
		c.err = err
	}
}

// implement the io.Reader interface
func (c *Reader) Read(p []byte) (n int, err error) {
	// checks for when we have no data
	for c.writePos == 0 || c.readPos == c.writePos {
		// if we have an error / EOF, just return it
		if c.err != nil {
			return n, c.err
		}
		// else, fill our buffer
		c.fillBuffer()
	}
	// we should have an appropriate amount of data, convert it into the given buffer
	bytesRead, bytesWritten, err := c.converter.Convert(c.buffer[c.readPos:c.writePos], p)
	// adjust byte counters
	c.readPos += bytesRead
	n += bytesWritten
	// if we experienced an iconv error, check it
	if err != nil {
		// E2BIG errors can be ignored (we'll get them often) as long
		// as at least 1 byte was written. If we experienced an E2BIG
		// and no bytes were written then the buffer is too small for
		// even the next character
		if err != syscall.E2BIG || bytesWritten == 0 {
			// track anything else
			c.err = err
		}
	}
	// return our results
	return n, c.err
}
