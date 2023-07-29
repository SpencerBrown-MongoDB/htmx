package main

import (
	"log"

	"github.com/SpencerBrown-MongoDB/htmx/contact"
	"github.com/SpencerBrown-MongoDB/htmx/web"
)

func main() {
	err := contact.GetContacts()
if err != nil {
	log.Fatalf("error loading contacts: %v", err)
}
	web.Setup()
}
