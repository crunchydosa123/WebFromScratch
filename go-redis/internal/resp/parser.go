package resp

import (
	"bufio"
	"errors"
	"strconv"
)

func ParseArray(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if line[0] != '*' {
		return nil, errors.New("Expected array")
	}

	n, _ := strconv.Atoi(line[1 : len(line)-2])
	result := make([]string, 0, n)

	for i := 0; i < n; i++ {
		lenLine, _ := r.ReadString('\n')
		size, _ := strconv.Atoi(lenLine[1 : len(lenLine)-2])

		buf := make([]byte, size+2)
		r.Read(buf)

		result = append(result, string(buf[:size]))
	}

	return result, nil
}
