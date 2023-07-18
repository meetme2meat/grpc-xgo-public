package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"xgo/main/src/gen"
	model "xgo/main/src/models"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddr = "0.0.0.0:8100"
)

type jwt struct {
	token string
}

func (j jwt) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j jwt) RequireTransportSecurity() bool {
	return true
}

type action func(context.Context, gen.CompanyServiceClient)

var actionMapper = map[string]action{
	"create": create,
	"update": update,
	"delete": delete,
	"get":    get,
}

func main() {
	action := flag.String("action", "create", "usage get, create, update, delete")
	help := flag.Bool("help", false, "show this message")
	// token := flag.String("token", "", "Path to JWT auth token.")

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()
	// if *token == "" {
	// 	fmt.Println("please provide a jwt token")
	// 	os.Exit(1)
	// }

	// jwtCreds, err := NewFromTokenFile(*token)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ctx := context.Background()
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	// conn, err := grpc.Dial(serverAddr, opts, grpc.WithPerRPCCredentials(jwtCreds))
	conn, err := grpc.Dial(serverAddr, opts)

	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := gen.NewCompanyServiceClient(conn)
	if callback, ok := actionMapper[*action]; ok {
		callback(ctx, client)
	} else {
		fmt.Println("not a valid callback")
	}
}

func NewFromTokenFile(token string) (credentials.PerRPCCredentials, error) {
	data, err := ioutil.ReadFile(token)
	if err != nil {
		return jwt{}, err
	}
	return jwt{string(data)}, nil
}

func create(ctx context.Context, client gen.CompanyServiceClient) {
	log.Println("Saving test company data via company service")
	c := &gen.Company{
		Id:          uuid.New().String(),
		Name:        "Google",
		Description: nil,
		Employee:    10,
		Registered:  true,
		Type:        0,
	}

	resp, err := client.CreateCompany(ctx, &gen.CreateCompanyRequest{Company: c})
	if err != nil {
		log.Fatalf("create company: %v", err)
	}

	fmt.Println("company created", resp.Id)
}

func get(ctx context.Context, client gen.CompanyServiceClient) {
	fmt.Println("provide company id")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		panic(err)
	}

	id := scanner.Text()
	fmt.Println("querying for company Id", id)

	result, err := client.GetCompany(ctx, &gen.GetCompanyRequest{Id: id})
	if err != nil {
		fmt.Println("received error", err)
		return
	}

	c := model.CompanyFromProto(result.Company)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(&c)
	if err != nil {
		panic(err)
	}
}

func delete(ctx context.Context, client gen.CompanyServiceClient) {
	fmt.Println("provide company id")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		panic(err)
	}

	id := scanner.Text()
	fmt.Println("delete for company Id", id)

	result, err := client.DeleteCompany(ctx, &gen.DeleteCompanyRequest{Id: id})
	if err != nil {
		panic(err)
	}

	fmt.Println("company id delete", result.Id)
}

func update(ctx context.Context, client gen.CompanyServiceClient) {
	fmt.Println("please provide company Id")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		panic(err)
	}

	id := scanner.Text()

	desc := "One of the largest company"
	c := &gen.Company{
		Id:          id,
		Name:        "Microsoft",
		Description: &desc,
		Employee:    20,
		Registered:  true,
		Type:        1,
	}

	fmt.Println("patch for company Id", id)
	result, err := client.PatchCompany(ctx, &gen.PatchCompanyRequest{Id: id, Company: c})
	if err != nil {
		panic(err)
	}

	newC := model.CompanyFromProto(result.Company)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(&newC)
	if err != nil {
		panic(err)
	}
}
