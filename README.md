# Goncurrent

Golang concurrency library for noobs like a me (Vay3t) [WIP]

# How to use

```go
package main

import (
	"fmt"
	"time"

	"github.com/nchgroup/goncurrent"
)

func main() {
	threads := 5

	// create a custom list with whatever type you want
	list := []string{"hello", "world", "from", "go", "goroutines"}

	goncurrent.Execute(threads, list, func(item interface{}) { // don't edit this line

		// you can write whatever
		fmt.Println(item)
		time.Sleep(1 * time.Second)
		// end whatever

	})
}
```

# Authors

* Vay3t - https://twitter.com/vay3t - https://gitlab.com/vay3t
* Gato - https://www.linkedin.com/in/daniel-mena-342a4715a/
