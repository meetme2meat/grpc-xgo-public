== README

```
go mod tidy
```


```
git clone github.com/bufbuild/protoc-gen-validate
cd protoc-gen-validate && make build
```

Generate gprc client and server
```
protoc company.proto \
    -I=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.0.2 \
    --go_opt=paths=source_relative \
    --go_out=src/gen \
    --go-grpc_opt=paths=source_relative \
    --go-grpc_out=src/gen \
    --proto_path=./src/api \
    --validate_out=paths=source_relative,lang=go:src/gen
```

Generate mock 

```
mockgen -source=src/internal/controller/company/controller.go -destination=src/internal/controller/company/mock/repository.go
```


```
docker compose -f docker-compose.yml up
```

a sample grpc client is provided 
```
    go run ./utility/client/client.go  -h

```

> Create a company
```
   go run ./main/src/utility/client/client.go -action create

```

> 2023/07/17 16:28:32 Saving test company data via company service
> company created d4edfccf-2561-43b0-8ca8-8259c2436eaf


> Update Company
```

go run ./main/src/utility/client/client.go -action update
please provide company Id
d4edfccf-2561-43b0-8ca8-8259c2436eaf
```

Result
```
patch for company Id d4edfccf-2561-43b0-8ca8-8259c2436eaf
{
    "id": "d4edfccf-2561-43b0-8ca8-8259c2436eaf",
    "name": "Microsoft",
    "description": "One of the largest company",
    "employee_count": 20,
    "registered": true,
    "type": "NONPROFIT"
}
```

> Get Company
```
go run ./main/src/utility/client/client.go -action get
provide company id
d4edfccf-2561-43b0-8ca8-8259c2436eaf
```

querying for company Id d4edfccf-2561-43b0-8ca8-8259c2436eaf
{
    "id": "d4edfccf-2561-43b0-8ca8-8259c2436eaf",
    "name": "Microsoft",
    "description": "One of the largest company",
    "employee_count": 20,
    "registered": true,
    "type": "NONPROFIT"
}
```



> Kafka ISSUE:

Kafka was complaining a lot when attempting to build it on alpine. For now event publisher/producer Nats is used. I will be raising a ticket on kafka client github page.


> JWT token:
  Not working because of some TLS CA cert issue.

== NOTE

Hello Fellow Codereviewer. I hope you read this, my intention of writing this is to have your feedback for my work.

I'm sure you would a lot of works so even if immediate feedback(I'm attaching my email id) is not possible
  please write me back about my work like..
    - what was good?
    - what went wrong?
    - how I can improve?
    - anything in general ...

  Thank you.


