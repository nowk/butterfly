# butterfly

[![Build Status](https://travis-ci.org/nowk/butterfly.svg?branch=master)](https://travis-ci.org/nowk/butterfly)
[![GoDoc](https://godoc.org/github.com/nowk/butterfly?status.svg)](http://godoc.org/github.com/nowk/butterfly)

Transform through pipe

## Example

    in, _ := os.Open("./words.txt")
    out, _ := os.Create("transformed.txt")

    tr := butterfly.NewTransform(in)
    tr.Through(UpperCase)
    tr.Through(Reversed)
    n, err := tr.WriteTo(out)
    if err != nil {
      // handle error
    }

    log.Printf("Wrote %d bytes", n)

## License

MIT