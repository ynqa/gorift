package errors

import (
	"bytes"
)

type MergedError []error

func (e MergedError) Error() string {
	var buf bytes.Buffer
	for _, err := range e {
		buf.WriteString(err.Error() + ": ")
	}
	res := buf.String()
	if buf.Len() >= 2 {
		res = res[:buf.Len()-2]
	}
	return res
}

func (e MergedError) Len() int {
	return len(e)
}

func (e MergedError) Add(err error) {
	if merged, ok := err.(MergedError); ok {
		e = append(e, merged...)
	} else {
		e = append(e, err)
	}
}
