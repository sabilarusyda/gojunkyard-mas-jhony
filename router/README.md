# Router
Package router implements http request multiplexer.

## Features
1. Has two engines. (HTTP router and Muxie router)
2. Support router grouping
3. Support middleware grouping
4. Can use standard mux middleware
5. Support parameterized url. `/entity/:id`

## What engine should I use?
There are several pros and cons about those engine

### HTTP Router
[https://github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)

##### Pros
* Recommended by many developer due to stability and performance

##### Cons
* Cannot handle conflicting route. Example: `/users/:id` and `/users/blocked`

### Muxie
[https://github.com/kataras/muxie](https://github.com/kataras/muxie)

##### Pros
* Can handle conflicting route.

##### Cons
* Not many review for this router. But the developer claims that the performance better than httprouter

## How to use?
```
package main

import (
	"fmt"
	"net/http"

	"devcode.xeemore.com/systech/gojunkyard/router"
)

func main() {
	r := router.NewWithEngine(router.Muxie)
	ruser := r.Group("/users" /*, middleware... */)
	ruser.GET("/:id", func(w http.ResponseWriter, r *http.Request) {
		id := router.GetParam(r, "id")
		fmt.Fprintln(w, "hoho", id)
	})
	ruser.GET("/ets", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":8080", r)
}
```
