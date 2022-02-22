# Pipeliner
Pipeliner is used for batching the request based on window and limit

## How to use
### Without Context
```
package main

import (
	"fmt"
	"time"

	"devcode.xeemore.com/systech/gojunkyard/pipeliner"
)

func f(v []int) error {
	// do something
	return nil
}

func main() {
	pipeliner := pipeliner.New(
		f,
		pipeliner.SetConcurrency(10),
		pipeliner.SetWindow(time.Second/2, 10),
	)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			pipeliner.Do(i)
		}(i)
		time.Sleep(time.Millisecond * 200)
	}

	time.Sleep(time.Minute)
}
```

### With Context
```
package main

import (
	"fmt"
	"time"

	"devcode.xeemore.com/systech/gojunkyard/pipeliner"
)

func f(ctx context.Context, v []int) error {
	// do something
	return nil
}

func main() {
	pipeliner := pipeliner.New(
		f,
		pipeliner.SetConcurrency(10),
		pipeliner.SetWindow(time.Second/2, 10),
	)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			pipeliner.Do(i)
		}(i)
		time.Sleep(time.Millisecond * 200)
	}

	time.Sleep(time.Minute)
}
```