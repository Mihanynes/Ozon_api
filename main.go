package main

import (
	"context"
	"fmt"
	"github.com/Mihanynes/Ozon_api.git/convert_functions"
	"github.com/Mihanynes/Ozon_api.git/converter"
	"github.com/Mihanynes/Ozon_api.git/storage"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type myConverter struct {
	storage storage.KeyValueStorage
	converter.UnimplementedConverterServer
}

func (s myConverter) PostUrl(con context.Context, req *converter.CreateRequest) (*converter.CreateResponse, error) {
	key := req.Str

	str, err := convert_functions.Convert_url_to_string(req.GetStr())
	if err != nil {
		return nil, err
	}

	// Store the converted URL string with a unique key
	err = s.storage.Set(key, str)
	if err != nil {
		return nil, err
	}

	return &converter.CreateResponse{
		Str: str,
	}, nil
}

func (s myConverter) GetUrl(con context.Context, req *converter.CreateRequest) (*converter.CreateResponse, error) {
	str, err := s.storage.GetKeyByValue(req.GetStr())
	if err != nil {
		return nil, err
	}
	if len(str) != convert_functions.ShortURLLength {
		fmt.Errorf("Unexpected length of converted string. Expected: %d, Got: %d", convert_functions.ShortURLLength, len(str))

	}

	for _, char := range str {
		if !convert_functions.Contains([]byte(convert_functions.AllowedChars), byte(char)) {
			fmt.Errorf("Unexpected character in converted string: %c", char)
		}
	}

	return &converter.CreateResponse{
		Str: str,
	}, nil
}

func main() {
	//var param int
	fmt.Println("Enter parametr of database:\n1 - PostreSQL\n2 - InMemory")
	//fmt.Scanf("%d\n", &param)
	param := os.Getenv("MY_PARAMETER")

	var service *myConverter
	switch param {
	case "2":
		inMemoryStorage := storage.NewInMemoryKeyValueStorage()
		storage.UseKeyValueStorage(inMemoryStorage)
		service = &myConverter{storage: inMemoryStorage}
	case "1":
		postgreSQLStorage, err := storage.NewPostgreSQLKeyValueStorage()
		if err != nil {
			log.Fatal(err)
		}
		storage.UseKeyValueStorage(postgreSQLStorage)
		service = &myConverter{storage: postgreSQLStorage}
	default:
		fmt.Println("Incorrect parameter")
		return

	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	} else {
		fmt.Println("Listening to server")
	}

	serverRegistrar := grpc.NewServer()
	//service := &myConverter{}
	converter.RegisterConverterServer(serverRegistrar, service)

	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}

	return
}
