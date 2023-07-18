package testutil

import (
	"log"
	"xgo/main/src/event"
	gen "xgo/main/src/gen"
	companyCtrl "xgo/main/src/internal/controller/company"
	grpchandler "xgo/main/src/internal/handler/grpc"
	"xgo/main/src/internal/repository"

	"github.com/golang-jwt/jwt"
	"github.com/gookit/config/v2"
)

var publicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAz2IZkElE8rmZtk3iLBVR
WxIagH6GEDlVnVjmInzm009Uu9fYhNeq7p8bLFyliaDUUDouE1d6sKT/KDdiieCA
JrV87HDc2FsOJ0CeT/aBFjBh8vM4asPJdZz+Ob/QTBDe8q37oAY/VUAnoZw6PzeS
Ic6xZn51aele/9U7E8d+GCHcGX3PIuTHljpcPPTktVNvij3dMZBuYXAc3EWQk5zV
Zx2pvMTmey9YGs/tHAjjxPkbK6JYqM3lrFEd/OQl+bb33QbCRVPSH3tGtaAjtged
Z78utLQa5J7sbluYS5DdwbIs1AJ33pFHeZLztUUqXYISoQV6/bq1gaa7clfU1YfE
/7XZBiRl3BbvqtQPejTA5NmKyZ9fwhomvdM+1oDXt0H9jvEREA0W8/rZVuC5v5Yc
oWgll/2ov+XPAQimwQibPOSjH+6B2TTZ9aNvAezjkHK5W0fr/WSQJKzEEqVnZkOJ
SQl2r+Zc6zbaqEIODVpIWJoEYz3+c5A5XQQNgWIhMxLIK6WYyxKXSrvQOMh408T3
YNJDUP8+Oolwa/xtYghlb7bib/r2BFoDrLbonijnIBYhQx8rbR07Atowvcb+JLJC
hbn86mUF9iqvBZlJqf3j3JJVwzkZQbDLcsMHGM3aatiXmRD+TGV/OLhzw/KAerXH
m/OJrgR3uGaS89JtkT3TqckCAwEAAQ==
-----END PUBLIC KEY-----`

func NewTestCompanyGRPCServer(ch chan event.Event) gen.CompanyServiceServer {
	config.Set("main.uri", "test-gorm.db")
	repo := repository.New("sqlite", ch)
	ctrl := companyCtrl.New(repo)

	publickey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		log.Panicf("Error parsing the jwt public key: %s", err)
	}

	return grpchandler.New(ctrl, publickey)
}
