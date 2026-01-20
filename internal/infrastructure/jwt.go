package infrastructure

import (
	"github.com/HadesHo3820/ebvn-golang-course/pkg/common"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils"
)

// CreateJWTProvider creates jwtutils.JWTGenerator and jwtutils.JWTValidator
func CreateJWTProvider() (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGen, err := jwtutils.NewJWTGenerator("./private.pem")
	common.HandleError(err)
	jwtValidator, err := jwtutils.NewJWTValidator("./public.pem")
	common.HandleError(err)
	return jwtGen, jwtValidator
}
