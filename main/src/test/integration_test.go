package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"

	genAuth "xgo/auth/src/gen"
	authtest "xgo/auth/src/pkg/testutil"
	"xgo/main/src/event"
	gen "xgo/main/src/gen"

	companytest "xgo/main/src/pkg/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serviceName = "company"
	serverAddr  = "localhost:8082"
	authAddr    = "localhost:8900"
)

var serverKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAw5iwaF5PADO2BjzU4e3l9UmH8XbSUVosRMcYIUOwq2uL3ZyR
RvY0laLhu2hP4f3BUlIMM17DE830s+HWzMkBNE8hkZ91XCfR4NnpTFSFIO5EfeUi
w2I+SGQu9xCJNqwL/TSCkBEAiVsfcjZqUV26zY7TVdF+/xk31mGibT7vA8POGXOP
5QdMAV5NYLEFgRSdfQ62q6XTsrY6vgGnlGDCvuQgr9xXoUejdiltbxtwqG/arMEu
Z6rOz0XnRePMLthieV2QxLJY/tqZZ4Mf8tfueNlE3aLiP7DhRFjhnZlf9DcZiroG
om7zKrCmz0Lh9L8RR8yv5E1cIRkoAktiLIvxrwIDAQABAoIBAAzKqY7JzCTZPOg/
hjSYWFeoTWmvOaX0XbzJwHw8bwtm5yjBGocnhtzaYCTfd4nyDHiTwRSC+AMNjxlM
hb2yz49aNXnOkeBLAmDQH3/Mb0BuFLCfEZzxid02IBQsUqzup4IRsxA07HIPMYlI
ob9cf+D9nDiakNGiFpLAo9y9Juh6Alf012X3HPRtijLXgl76fR8mzxyOD1ZGaRSa
nm9Agnv6v8Zzx2cQXbfhEzjld/b4NCiR3JedYs1Mehj2dshusjK7Rk+EhYovql4u
akABCnvVBeBlUmoKuWYDfg6n+RjdhoHR+qjDhch4dLwnAjA6v4gnRQbrMhd97f+1
RGkcbZECgYEA9KRXKDUynVYCN6klClIOYNBaspRsbFBqj/PeLmi9AOLFJ4gt843m
YpHe5pmkLGYlwJjILDI85CXn7H/fg7Z+GQzhfI38DOyNDD3t9pS7lpeqOuwOYTnP
tppBGhQDDANuSxyXqqeYSOxEGyRTJ32xky8YzANgYvrwMAk+QMaV5W0CgYEAzK1s
uEX1t6f8PdUSD5T9oUfDxNDLB8wIF45Oxlk3bImAWMUI8wR2HK5eK+1vkeGEkbq8
X6gca8Fr33Nv7jmrANO6fl18iRKW3a/OtwVGmj8blaAoOd0TfgClsTCWs3Ubgl/F
JyJ86zrTcmzz45q3BMrEizqgOGvNnPeKyz73rgsCgYEAqY0grsg33R2YCdWby0xV
lLmysmP1xRfy0vQUf5utqmiAdcaG+m7VRmmMz8uaIf9lmNcKnL7wvrqaw6lYUuPu
/xOTT4zkLFzh4KMnQqeQX22b2Jxz1uSHVioQhq9p8TCLh1k4sFjZTWkaRqllTFBr
+vNAP1zzt4XtY410bNZ1Wv0CgYBdwHZTNeBmXnDg1a8vKfy/GkMm7MiC6scuGwYk
PotvkNAUWTRPNFTxsED8eAap2JXDtrhATJ2wEenacWLsyMd2WoVLCoFXvAcUxkm2
dZkwYAW/lJu4XXZnOd6reekdjF+saTfCRD7Z9JkUCanxMFXywPokGBd5oI+O/ag6
jr4enwKBgE+45hul9mVAMAj2EnWe96U9EgeWUlxvLSvyCwPx7GPMTeQU0KxuScrr
ny/cdvmE8yvuQuv8jmgUI80pkkU/Lr4mP3x+eRnRi0C7FseHxT5o5n8AXARUTfEu
15G1FpUdm/F5fjKe8E75X8WwcV0L5a7CQfpwpbErlG2182pkUXZD
-----END RSA PRIVATE KEY-----`

var serverCert = `-----BEGIN CERTIFICATE-----
MIIDuTCCAqGgAwIBAgIJALlmB71alZ+2MA0GCSqGSIb3DQEBCwUAMDkxEjAQBgNV
BAMMCWRlbW8ueGNvbTELMAkGA1UEBhMCVVMxFjAUBgNVBAcMDVNhbiBGcmFuc2lz
Y28wHhcNMjMwNzE3MTAwMjA1WhcNMjQwNzE2MTAwMjA1WjB2MQswCQYDVQQGEwJV
UzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5zaXNjbzEO
MAwGA1UECgwFVmlyZW4xFjAUBgNVBAsMDVZpcmVuZHJhIE5lZ2kxEjAQBgNVBAMM
CWRlbW8ueGNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMOYsGhe
TwAztgY81OHt5fVJh/F20lFaLETHGCFDsKtri92ckUb2NJWi4btoT+H9wVJSDDNe
wxPN9LPh1szJATRPIZGfdVwn0eDZ6UxUhSDuRH3lIsNiPkhkLvcQiTasC/00gpAR
AIlbH3I2alFdus2O01XRfv8ZN9Zhom0+7wPDzhlzj+UHTAFeTWCxBYEUnX0Otqul
07K2Or4Bp5Rgwr7kIK/cV6FHo3YpbW8bcKhv2qzBLmeqzs9F50XjzC7YYnldkMSy
WP7amWeDH/LX7njZRN2i4j+w4URY4Z2ZX/Q3GYq6BqJu8yqwps9C4fS/EUfMr+RN
XCEZKAJLYiyL8a8CAwEAAaOBhjCBgzBTBgNVHSMETDBKoT2kOzA5MRIwEAYDVQQD
DAlkZW1vLnhjb20xCzAJBgNVBAYTAlVTMRYwFAYDVQQHDA1TYW4gRnJhbnNpc2Nv
ggkAxXWaqxat6VswCQYDVR0TBAIwADALBgNVHQ8EBAMCBPAwFAYDVR0RBA0wC4IJ
ZGVtby54Y29tMA0GCSqGSIb3DQEBCwUAA4IBAQA6f8Q3vBkwhgKjzBic9o9sDU12
LilNtSLMJuYqNoX2C6sVda/G7v1luxm1BVjw1lYEyryHh+6rO69OX05xRUg/IrG4
akOeOFdM/VkXpjcdovItIUe4Eo8XxzVOFGmkDqgT1NYSt6T2dSSsZ+LjHU/fodvU
mobkU28T+dXFP3Zi/EVdUdwEqHKxk5lHgCcKhXlZf+VUssef2aiTFaB+TY8//mwE
zg9fM5wPQWkjlfvPU4oUPzx5yGAx7KB/RbE+2zpZfVJftSZK++ow6PIGlZjLq/++
fSdTxMfo/F3G+YvkCenxEladNTFwUibaPZxCVbu7IKyUpFUtTyPAzWcqMhMZ
-----END CERTIFICATE-----`
var jwtToken credentials.PerRPCCredentials

type jwt struct {
	token string
}

func NewJwtToken(token string) credentials.PerRPCCredentials {
	return jwt{string(token)}
}

func (j jwt) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j jwt) RequireTransportSecurity() bool {
	return true
}

func TestMain(m *testing.M) {

	srv := startAuthServer(context.TODO())
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(authAddr, opts)
	if err != nil {
		panic(err)
	}

	client := genAuth.NewAuthClient(conn)

	signReq := genAuth.SignupRequest{
		Username:        "admin",
		Password:        "admin123",
		PasswordConfirm: "admin123"}
	_, err = client.Signup(context.TODO(), &signReq)
	if err != nil {
		panic(err)
	}

	loginReq := genAuth.LoginRequest{
		Username: "admin",
		Password: "admin123"}

	tokenResp, err := client.Login(context.TODO(), &loginReq)
	if err != nil {
		panic(err)
	}

	fmt.Println("received token")
	jwtToken = NewJwtToken(tokenResp.GetToken())
	m.Run()
	conn.Close()
	srv.GracefulStop()
}
func TestIntegration(t *testing.T) {
	log.Println("Starting the integration test")

	ctx := context.Background()

	log.Println("Setting up service handlers and clients")
	ch := make(chan event.Event, 1)
	companySrv := startCompanyServer(ctx, ch)
	defer companySrv.GracefulStop()

	// creds := tls.GetTransportCredentials("./server.crt", "./server.key")
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	// // Please go through README.md
	// grpc.WithPerRPCCredentials(jwtToken))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := gen.NewCompanyServiceClient(conn)

	log.Println("Saving test company data via company service")

	c := &gen.Company{
		Id:          uuid.New().String(),
		Name:        "Google",
		Description: nil,
		Employee:    10,
		Registered:  true,
		Type:        0,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var data event.Event
	go func() {
		defer wg.Done()
		data = <-ch
	}()
	resp, err := client.CreateCompany(ctx, &gen.CreateCompanyRequest{Company: c})
	if err != nil {
		log.Fatalf("create company: %v", err)
	}

	assert.Equal(t, c.Id, resp.Id, "must be equal")

	log.Println("Retrieving company ...")

	result, err := client.GetCompany(ctx, &gen.GetCompanyRequest{Id: resp.Id})
	if err != nil {
		log.Fatalf("get company: %v", err)
	}
	c.Id = resp.Id
	assert.Equal(t, result.Company.Description, c.Description, "must be equal")
	assert.Equal(t, result.Company.Id, c.Id, "must be equal")
	assert.Equal(t, result.Company.Employee, c.Employee, "must be equal")
	assert.Equal(t, result.Company.Registered, c.Registered, "must be equal")
	assert.Equal(t, result.Company.Type, c.Type, "must be equal")

	wg.Wait()
	assert.Equal(t, data.RecordID, c.Id)
	assert.Equal(t, data.RecordType, "companies")
	assert.Equal(t, data.EventType, "create")
	log.Println("Update Integration test execution successful")

	c.Type = 2
	description := "internet search company"
	c.Description = &description
	c.Employee = 12
	wg.Add(1)
	go func() {
		defer wg.Done()
		data = <-ch
	}()

	patchResult, err := client.PatchCompany(ctx, &gen.PatchCompanyRequest{Company: c, Id: c.Id})
	if err != nil {
		log.Fatalf("patch company: %v", err)
	}

	assert.Equal(t, *patchResult.Company.Description, "internet search company", "must be equal")
	assert.Equal(t, patchResult.Company.Id, c.Id, "must be equal")
	assert.Equal(t, patchResult.Company.Employee, uint32(12), "must be equal")
	assert.Equal(t, patchResult.Company.Registered, c.Registered, "must be equal")
	assert.Equal(t, patchResult.Company.Type, gen.Type(2), "must be equal")
	wg.Wait()
	assert.Equal(t, data.RecordID, c.Id)
	assert.Equal(t, data.RecordType, "companies")
	assert.Equal(t, data.EventType, "update")

	wg.Add(1)
	go func() {
		defer wg.Done()
		data = <-ch
	}()

	deleteRslt, err := client.DeleteCompany(ctx, &gen.DeleteCompanyRequest{Id: c.Id})
	if err != nil {
		log.Fatalf("delete company: %v", err)
	}

	assert.Equal(t, deleteRslt.Id, c.Id, "must be equal")
	wg.Wait()
	assert.Equal(t, data.RecordID, c.Id, "event record id should match")
	assert.Equal(t, data.RecordType, "companies")
	assert.Equal(t, data.EventType, "delete")

	result, err = client.GetCompany(ctx, &gen.GetCompanyRequest{Id: resp.Id})
	assert.Nil(t, result, "must be nil")
	assert.ErrorContains(t, err, "company data not found")
}

func startCompanyServer(_ context.Context, ch chan event.Event) *grpc.Server {
	log.Println("Starting company service on " + serverAddr)
	h := companytest.NewTestCompanyGRPCServer(ch)
	l, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// // Please go through README.md

	// creds := tls.GetTransportCredentials("./server.crt", "./server.key")
	// srv := grpc.NewServer(grpc.Creds(creds))
	srv := grpc.NewServer()
	gen.RegisterCompanyServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()

	return srv
}

func startAuthServer(_ context.Context) *grpc.Server {
	log.Println("Starting auth server on " + authAddr)
	h := authtest.NewTestAuthServer()
	l, err := net.Listen("tcp", authAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	genAuth.RegisterAuthServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()

	return srv
}
