# butterfly

[![Build Status](https://travis-ci.org/nowk/butterfly.svg?branch=master)](https://travis-ci.org/nowk/butterfly)

Transform through pipe

## Example

    in, _ := os.Open("./words.txt")
    out, _ := os.Create("transformed.txt")

    tr := butterfly.NewCocoon(in)
    tr.Through(UpperCase)
    tr.Through(Reversed)
    n, err := tr.Write(out)
    if err != nil {
      // handle error
    }

    log.Printf("Wrote %d bytes", n)

## License

MIT