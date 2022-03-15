package iconv

import "io"

type Writer struct {
	destination       io.Writer
	converter         *Converter
	buffer            []byte
	readPos, writePos int
	err               error
}

func NewWriter(destination io.Writer, fromEncoding string, toEncoding string) (*Writer, error) {
	// create a converter
	converter, err := NewConverter(fromEncoding, toEncoding)
	if err == nil {
		return NewWriterFromConverter(destination, converter), err
	}
	// return the error
	return nil, err
}

func NewWriterFromConverter(destination io.Writer, converter *Converter) (writer *Writer) {
	writer = new(Writer)
	// copy elements
	writer.destination = destination
	writer.converter = converter
	// create 8K buffers
	writer.buffer = make([]byte, 8*1024)
	return writer
}

func (c *Writer) emptyBuffer() {
	// write new data out of buffer
	bytesWritten, err := c.destination.Write(c.buffer[c.readPos:c.writePos])
	// update read position
	c.readPos += bytesWritten
	// slide existing data to beginning
	if c.readPos > 0 {
		// copy current bytes - is this guaranteed safe?
		copy(c.buffer, c.buffer[c.readPos:c.writePos])
		// adjust positions
		c.writePos -= c.readPos
		c.readPos = 0
	}
	// track any reader error / EOF
	if err != nil {
		c.err = err
	}
}

// implement the io.Writer interface
func (c *Writer) Write(p []byte) (n int, err error) {
	// write data into our internal buffer
	bytesRead, bytesWritten, err := c.converter.Convert(p, c.buffer[c.writePos:])
	// update bytes written for return
	n += bytesRead
	c.writePos += bytesWritten
	// checks for when we have a full buffer
	for c.writePos > 0 {
		// if we have an error, just return it
		if c.err != nil {
			return
		}
		// else empty the buffer
		c.emptyBuffer()
	}
	return n, err
}
