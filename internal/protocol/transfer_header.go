package protocol

import (
	"bufio"
	"fmt"
)

func writeTransferHeader(w *bufio.Writer, t Transfer) (int, error) {
	return fmt.Fprintf(w, "$ORIGIN %s\n$SERIAL %d\n", t.Zone, t.Serial)
}
