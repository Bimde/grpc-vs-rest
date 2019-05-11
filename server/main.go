package main

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/Bimde/grpc-vs-rest/pb"
)

func handle(w http.ResponseWriter, _ *http.Request) {
	random := pb.Random{RandomString: "a_random_string", RandomInt: 1984}
	bytes, err := json.Marshal(&random)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func main() {
	server := &http.Server{Addr: "bimde:8080", Handler: http.HandlerFunc(handle)}
    log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}