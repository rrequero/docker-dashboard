package main

import (
	"fmt"

	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {
	do.Env = `GOPATH=$GOPATH`

	p.Task("default", do.S{"build"}, nil)

	p.Task("build", nil, func(c *do.Context) {
		c.Run("GOOS=linux GOARCH=amd64 go build", do.M{"$in": ""})
	}).Src("**/*.go")

	p.Task("server", nil, func(c *do.Context) {
		// rebuilds and restarts when a watched file changes
		c.Start("main.go", do.M{"$in": ""})
	}).Src("**/*.go").Debounce(100)

	p.Task("test", nil, func(c *do.Context) {
		if c.FileEvent != nil {
			fmt.Println(c.FileEvent.Path) // => /path/to/this/file.go
		}
	}).Src("**/*.go")

}

func main() {
	do.Godo(tasks)
}
