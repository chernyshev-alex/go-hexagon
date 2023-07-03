//go:generate wire
package cmd

import (
	"log"
)

func main() {
	if app, err := initializeWire(); err != nil {
		log.Fatal(err)
	} else {
		app.Listen(":3000")
	}
}
