package main

import (
	"ex3/cyoa"
	"fmt"
	"log"
	"net/http"
)

func main() {
	story := cyoa.LoadJson("../gopher.json")

	handler := cyoa.NewHandler(story, nil)

	fmt.Printf("Started server at Port %d\n", 8080)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), handler))

	fmt.Print(story)
}
