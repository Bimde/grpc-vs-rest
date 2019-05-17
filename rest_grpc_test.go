package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/Bimde/grpc-vs-rest/pb"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
)

func BenchmarkHTTP2GetWithWokers(b *testing.B) {
	client.Transport = &http2.Transport{
		TLSClientConfig: createTLSConfigWithCustomCert(),
	}
	requestQueue := make(chan Request)
	defer startWorkers(&requestQueue, noWorkers, startPostWorker)()
	b.ResetTimer() // don't count worker initialization time
	for i := 0; i < b.N; i++ {
		requestQueue <- Request{
			Path: "https://bimde:8080",
			Random: &pb.Random{
				RandomInt:    2019,
				RandomString: "a_string",
			},
		}
	}
}

func BenchmarkGRPCWithWokers(b *testing.B) {
	conn, err := grpc.Dial("bimde:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	client := pb.NewRandomServiceClient(conn)
	requestQueue := make(chan Request)
	defer startWorkers(&requestQueue, noWorkers, getStartGRPCWorkerFunction(client))()
	b.ResetTimer() // don't count worker initialization time

	for i := 0; i < b.N; i++ {
		requestQueue <- Request{
			Path: "http://localhost:9090",
			Random: &pb.Random{
				RandomInt:    2019,
				RandomString: "a_string",
			},
		}
	}
}

func post(path string, input interface{}, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		log.Println("error marshalling input ", err)
		return err
	}
	body := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		log.Println("error creating request ", err)
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("error executing request ", err)
		return err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading response body ", err)
		return err
	}

	err = json.Unmarshal(bytes, output)
	if err != nil {
		log.Println("error unmarshalling response ", err)
		return err
	}

	return nil
}

func getStartGRPCWorkerFunction(client pb.RandomServiceClient) func(*chan Request, *sync.WaitGroup) {
	return func(requestQueue *chan Request, wg *sync.WaitGroup) {
		go func() {
			for {
				request := <-*requestQueue
				if request.Path == stopRequestPath {
					wg.Done()
					return
				}
				client.DoSomething(context.TODO(), request.Random)
			}
		}()
	}
}

func startPostWorker(requestQueue *chan Request, wg *sync.WaitGroup) {
	go func() {
		for {
			request := <-*requestQueue
			if request.Path == stopRequestPath {
				wg.Done()
				return
			}
			post(request.Path, request.Random, request.Random)
		}
	}()
}
