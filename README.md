# go-tf-random

[![CI](https://github.com/susisu/go-tf-random/workflows/CI/badge.svg)](https://github.com/susisu/go-tf-random/actions?query=workflow%3ACI)

A splittable pseudorandom number generator using [the Threefish block cipher](https://en.wikipedia.org/wiki/Threefish).

- [The original Haskell implementation](https://hackage.haskell.org/package/tf-random)
- [My JavaScript/TypeScript implementation](https://github.com/susisu/tf-random.js)

This module is based on the JS/TS implementation, and there are some minor differences (including bug fixes) compared to the original one.

## Usage

Use `go get` to install:

``` shell
go get github.com/susisu/go-tf-random
```

go-tf-random exports a generator implementation `TFGen` that allows you to generate random uint32 values.
To generate random values of other numeric types, you can use [go-random](https://github.com/susisu/go-random).

``` go
package main

import (
	"fmt"

	random "github.com/susisu/go-random/uint32"
	tf_random "github.com/susisu/go-tf-random"
)

func main() {
	g1 := tf_random.NewTFGen(42, 42, 42, 42)
	v1 := random.Float64(g1)
	fmt.Printf("%f\n", v1)
}
```

And most importantly, `TFGen` can be split to create new generators that are independent of the original one.
These new generators can be safely passed around to other functions and goroutines, which allows you to make your code more predictable.

``` go
func main() {
	// ...

	g2 := g1.Split()
	v2 := random.Float64(g2)
	fmt.Printf("%f\n", v2)
}
```

## Performance

Sadly, `TFGen` isn't that fast due to its relatively complex computations.
Use it only in non-performance-critical code, or help me improve its performance!

```
goos: darwin
goarch: arm64
pkg: github.com/susisu/go-tf-random
BenchmarkTFGen_Uint32-10    	61052778	        18.85 ns/op	       8 B/op	       0 allocs/op
BenchmarkMathRand-10        	526618652	         2.341 ns/op	       0 B/op	       0 allocs/op
```

## License

[MIT License](http://opensource.org/licenses/mit-license.php)

## Author

Susisu ([GitHub](https://github.com/susisu), [Twitter](https://twitter.com/susisu2413))
