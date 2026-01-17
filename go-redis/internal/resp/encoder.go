package resp

import (
	"bufio"
	"strconv"
)

func WriteSimpleString(w *bufio.Writer, s string) error {
	_, err := w.WriteString("+" + s + "\r\n")
	return err
}

func WriteError(w *bufio.Writer, msg string) error {
	_, err := w.WriteString("-ERR" + msg + "\r\n")
	return err
}

func WriteBulkString(w *bufio.Writer, s *string) error {
	if s == nil {
		_, err := w.WriteString("$-1\r\n")
		return err
	}

	_, err := w.WriteString("$" + strconv.Itoa(len(*s)) + "\r\n" + *s + "\r\n")
	return err
}
