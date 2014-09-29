package butterfly

import "io"

// TransformFunc type represents the signature for a transform function
type TransformFunc func(io.Writer, io.Reader) error

// Transform calls f(w, r)
func (f TransformFunc) Transform(w io.Writer, r io.Reader) error {
	return f(w, r)
}

// Transform is a struct around a collection of transform functions. Which takes
// the Source and pipes them through each stage and writes out to the final
// destination
type Transform struct {
	stages   []TransformFunc
	Source   io.Reader
	bytesize int
}

// NewTransform returns a new Transform
func NewTransform(r io.Reader) *Transform {
	return NewTransformSize(r, 32*1024)
}

// NewTransformSize returns a new Transform whose final read size will be the given
// number
func NewTransformSize(r io.Reader, n int) *Transform {
	return &Transform{
		Source:   r,
		bytesize: n,
	}
}

// Through stacks TransformFuncs to be piped through on Write
func (c *Transform) Through(fn TransformFunc) *Transform {
	c.stages = append(c.stages, fn)
	return c
}

// WriteTo writes to the final output after being piped through the stages
func (c *Transform) WriteTo(out io.Writer) (int, error) {
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

	return writeto(out, r, c.bytesize)
}

// writeto is the last stage, this writes to the final source
func writeto(w io.Writer, r io.Reader, bytesize int) (int, error) {
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
