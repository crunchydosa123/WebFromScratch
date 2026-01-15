package resp

import (
	"bufio"
)

func WriteSimpleString(w *bufio.Writer, s string) error {
	_, err := w.WriteString("+" + s + "\r\n")
	return err
}

func WriteError(w *bufio.Writer, msg string) error {
	_, err := w.WriteString("-ERR" + msg + "\r\n")
	return err
}
