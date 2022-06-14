package main

import (
	"context"
	"gonum.org/v1/gonum/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	pb "nlo/nlo"
	"sync"
)

const (
	address = "localhost:50051"
)

var pool *sync.Pool

type Frequency struct {
	array []float64
}

func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			return new(Frequency)
		},
	}
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	c := pb.NewGetRespClient(conn)

	stream, err := c.Do(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("%v.Do(_) = _, %v", c, err)
	}
	initPool()
	p := pool.Get().(*Frequency)
	count := 0
	frequency := make([]float64, 1e5)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Do(_) = _, %v", c, err)
		}
		frequency[count] = data.GetFrequency()
		if count < 200 {
			log.Printf("Processed %d value, frequency %f", count, frequency[count])
		} else {
			p.array = frequency
			mean, std := stat.MeanStdDev(p.array, nil)
			meanErr := stat.StdErr(std, float64(len(p.array)))
			meanArr, stdArr := stat.MeanStdDev(frequency, nil)
			meanErrArr := stat.StdErr(stdArr, float64(len(frequency)))
			if count%200 == 0 {
				//log.Println(pool.Get())
				log.Printf("Frequency: %f. Processed %d value. Predicted values od mean: %.5f and STD: %.5f. Error: %.5f", frequency[count], count, mean, std, meanErr)
				log.Printf("Frequency: %f. Processed %d value. Predicted values od mean: %.5f and STD: %.5f. Error: %.5f", frequency[count], count, meanArr, stdArr, meanErrArr)
			}
		}
		count++
	}
}
