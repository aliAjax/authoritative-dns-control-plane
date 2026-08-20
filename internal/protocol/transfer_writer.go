package protocol

import (
	"bufio"
	"io"
)

func newTransferBuffer(w io.Writer) *bufio.Writer {
	return bufio.NewWriterSize(w, 4096)
}

func writeTransferPayload(w *bufio.Writer, t Transfer) error {
	if _, err := writeTransferHeader(w, t); err != nil {
		return err
	}
	return writeTransferRecords(w, t.Records)
}

func flushTransferBuffer(w *bufio.Writer) error {
	if w == nil {
		return nil
	}
	return w.Flush()
}

func finishTransfer(w *bufio.Writer, first error) error {
	flushErr := flushTransferBuffer(w)
	if first != nil {
		return first
	}
	return flushErr
}
