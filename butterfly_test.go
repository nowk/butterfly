package butterfly

import "fmt"
import "io"
import "bytes"
import "strings"
import "testing"
import "time"
import "github.com/nowk/assert"

func Hello(w io.Writer, r io.Reader) error {
	time.Sleep(500 * time.Millisecond)
	b := make([]byte, 6) // NOTE for test this matches exactly to "Dear, "
	n, _ := r.Read(b)
	b = append(b[:n], []byte("Hello")...)
	w.Write(b)
	return nil
}

func World(w io.Writer, r io.Reader) error {
	var c []byte
	for {
		time.Sleep(100 * time.Millisecond)
		b := make([]byte, 2)
		n, err := r.Read(b)
		if err == io.EOF || n == 0 {
			break
		}
		c = append(c, b[:n]...)
	}
	c = append(c, []byte(" World!")...)
	w.Write(c)
	return nil
}

func Uhoh(w io.Writer, r io.Reader) error {
	return fmt.Errorf("something broke")
}

func TestTransform(t *testing.T) {
	r := strings.NewReader("Dear, ")
	var b []byte
	buf := bytes.NewBuffer(b)

	tr := NewTransform(r)
	tr.Through(Hello)
	tr.Through(World)
	n, err := tr.WriteTo(buf)

	assert.Nil(t, err)
	assert.Equal(t, 18, n)
	assert.Equal(t, "Dear, Hello World!", buf.String())
}

// func TestError(t *testing.T) {
// 	r := strings.NewReader("Dear, ")
// 	var b []byte
// 	buf := bytes.NewBuffer(b)

// 	tr := NewTransform(r)
// 	tr.Through(Hello)
// 	tr.Through(Uhoh)
// 	tr.Through(World)
// 	n, err := tr.Write(buf)

// 	buf.Truncate(n)
// 	fmt.Println(buf.String())
// 	assert.NotNil(t, err)
// 	assert.Equal(t, "something broke", err.Error())
// 	assert.Equal(t, 0, n)
// }
