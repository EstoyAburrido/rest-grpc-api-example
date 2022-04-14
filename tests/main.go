package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/estoyaburrido/rest-grpc-api-example/tests/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type fibbonaciQuery struct {
	X uint64 `json:"x"`
	Y uint64 `json:"y"`
}

type fibbonaciIndex struct {
	MaxIndex uint64 `json:"maxIndex"`
}

var reference = []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}

const (
	firstHttpX  = 0
	firstHttpY  = 5
	secondHttpX = 2
	secondHttpY = 13
	firstGrpcX  = 0
	firstGrpcY  = 10
	secondGrpcX = 3
	secondGrpcY = 8
)

func init() {
	viper.SetDefault("Host", "web-api")
	viper.SetDefault("HttpPort", "8080")
	viper.SetDefault("GrpcPort", "9111")

	viper.SetConfigName("config/config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		exitErr(fmt.Errorf("fatal error with config file: %s", err))
	}
	viper.WatchConfig()
}

func main() {
	host := viper.GetString("Host")
	httpPort := viper.GetString("HttpPort")
	grpcPort := viper.GetString("GrpcPort")

	httpAddr := fmt.Sprintf("%v:%v", host, httpPort)
	grpcAddr := fmt.Sprintf("%v:%v", host, grpcPort)

	err := testRest(httpAddr)
	if err != nil {
		exitErr(err)
	}

	err = testGrpc(grpcAddr)
	if err != nil {
		exitErr(err)
	}

	log.Println("All tests passed")

	os.Exit(0)
}

func testRest(addr string) error {
	urlSeq := fmt.Sprintf("http://%v/getSequence", addr)
	urlIdx := fmt.Sprintf("http://%v/getMaxIndex", addr)

	log.Println("REST: first fibonacci sequence test")
	reqBody := fibbonaciQuery{
		X: firstHttpX,
		Y: firstHttpY,
	}
	body, err := post(urlSeq, reqBody)
	if err != nil {
		return err
	}

	var seqResp []uint64
	err = json.Unmarshal(body, &seqResp)
	if err != nil {
		return err
	}

	for i := range seqResp {
		if reference[i] != seqResp[i] {
			return fmt.Errorf("assertion failed at element #%v: expected %v but got %v", i, reference[i], seqResp[i])
		}
	}
	log.Println("OK")

	log.Println("REST: second fibonacci sequence test")
	reqBody.X = secondHttpX
	reqBody.Y = secondHttpY

	body, err = post(urlSeq, reqBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &seqResp)
	if err != nil {
		return err
	}

	newref := reference[secondHttpX : secondHttpY+1]
	for i := range seqResp {
		if newref[i] != seqResp[i] {
			return fmt.Errorf("assertion failed at element #%v: expected %v but got %v", i, reference[i], seqResp[i])
		}
	}
	log.Println("OK")

	log.Println("REST: max index test")
	body, err = post(urlIdx, nil)
	if err != nil {
		return err
	}
	idxResp := fibbonaciIndex{}
	err = json.Unmarshal(body, &idxResp)
	if err != nil {
		return err
	}

	if idxResp.MaxIndex != secondHttpY {
		return fmt.Errorf("assertion failed at max index check: expected %v but got %v", secondHttpY, idxResp)
	}
	log.Println("OK")

	return nil
}

func testGrpc(addr string) error {

	log.Println("GRPC: first fibonacci sequence test")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to GRPC: %v", err)
	}
	defer conn.Close()
	c := pb.NewFibonacciServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &pb.GetSequenceRequest{
		X: firstGrpcX,
		Y: firstGrpcY,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	res := r.GetRes()

	for i := range res {
		if reference[i] != res[i] {
			return fmt.Errorf("assertion failed at element #%v: expected %v but got %v", i, reference[i], res[i])
		}
	}
	log.Println("OK")

	log.Println("GRPC: second fibonacci sequence test")
	r, err = c.Get(ctx, &pb.GetSequenceRequest{
		X: secondGrpcX,
		Y: secondGrpcY,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	res = r.GetRes()
	newref := reference[secondGrpcX : secondGrpcY+1]

	for i := range res {
		if newref[i] != res[i] {
			return fmt.Errorf("assertion failed at element #%v: expected %v but got %v", i, reference[i], res[i])
		}
	}
	log.Println("OK")

	return nil
}

func post(url string, i interface{}) ([]byte, error) {
	bodyJson, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func exitErr(e error) {
	fmt.Println(e.Error())

	os.Exit(1)
}
