package butterfly

import "io"

// TransformFunc type represents the signature for a transform function
type TransformFunc func(io.Writer, io.Reader) error

// Transform calls f(w, r)
func (f TransformFunc) Transform(w io.Writer, r io.Reader) error {
	return f(w, r)
}

// Cocoon is a struct around a collection of transform functions. Which takes
// the Source and pipes them through each stage and writes out to the final
// destination
type Cocoon struct {
	stages   []TransformFunc
	Source   io.Reader
	bytesize int
}

// NewCocoon returns a new Cocoon
func NewCocoon(r io.Reader) *Cocoon {
	return NewCocoonSize(r, 32*1024)
}

// NewCocoonSize returns a new Cocoon whose final read size will be the given
// number
func NewCocoonSize(r io.Reader, n int) *Cocoon {
	return &Cocoon{
		Source:   r,
		bytesize: n,
	}
}

// Through stacks TransformFuncs to be piped through on Write
func (c *Cocoon) Through(fn TransformFunc) *Cocoon {
	c.stages = append(c.stages, fn)
	return c
}

// Write writes to the final output after being piped through the stages
func (c *Cocoon) Write(out io.Writer) (int, error) {
	r := c.Source

	for _, fn := range c.stages {
		p, w := io.Pipe()
		go func(fn TransformFunc, w io.Writer, r io.Reader) {
			err := fn.Transform(w, r)
			if err != nil {
				// TODO handle error
			}
			w.Write(nil)
		}(fn, w, r)

		r = p
	}

	return write(out, r, c.bytesize)
}

// write is the last stage, this writes to the final source
func write(w io.Writer, r io.Reader, bytesize int) (int, error) {
	i := 0

Write:
	for {
		b := make([]byte, bytesize)
		n, err := r.Read(b)
		if err == io.EOF || n == 0 {
			break Write
		}

		if err != nil {
			return i, err
		}

		m, err := w.Write(b[:n])
		i += m
		if err != nil {
			return i, err
		}
	}

	return i, nil
}
