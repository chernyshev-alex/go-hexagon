//go:generate wire
package cmd

import (
	"log"
)

func main() {
	app, err := initializeWire()
	if err != nil {
		log.Fatal(err)
	}
	app.Listen(":3000")
}
