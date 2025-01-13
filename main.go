package main

import "log"

func main() {
	store, err := CreateDb()
	if err != nil {
		log.Fatal(err)
	}

	err = store.init()
	if err != nil {
		log.Fatal(err)
	}
	api := MakeApi(":8080", store)
	api.run()
}
