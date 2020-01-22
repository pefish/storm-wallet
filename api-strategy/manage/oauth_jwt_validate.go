package manage

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/ory/hydra-client-go/client/public"
	"github.com/ory/hydra-client-go/models"
	"github.com/pefish/go-config"
	"github.com/pefish/go-reflect"
	"log"
	"math/big"
	"strings"
	"time"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/global"

	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
)

type OauthJwtValidateClass struct {
}

var OauthJwtValidateStrategy = OauthJwtValidateClass{}

func (this *OauthJwtValidateClass) GetName() string {
	return `oauth_jwt_validate`
}

func (this *OauthJwtValidateClass) GetDescription() string {
	return `校验授权服务器颁发的jwt`
}

func (this *OauthJwtValidateClass) GetErrorCode() uint64 {
	return constant.JWT_ERROR
}

type OauthJwtValidateParam struct {
	RequiredScopes []string
}

func (this *OauthJwtValidateClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {

	// 校验jwt合法性
	jwtStr := out.Ctx.GetHeader(`JSON-WEB-TOKEN`)
	if jwtStr == `` {
		go_error.ThrowInternal(`auth error. jwt not found.`)
	}

	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		cert, err := this.getPublicKey(token)
		if err != nil {
			panic(err)
		}
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, err
		}
		return verifyKey, nil
	})
	if err != nil {
		go_error.ThrowInternalError(`jwt is illegal`, err)
	}

	jwtBody := token.Claims.(jwt.MapClaims)
	out.JwtBody = jwtBody
	// 校验iss
	checkIss := jwtBody.VerifyIssuer(global.AuthServerUrl+`/`, false)
	if !checkIss {
		go_error.ThrowInternal(`Invalid issuer.`)
	}

	// 校验aud[0]必须是clientId
	checkAud := jwtBody.VerifyAudience(go_config.Config.GetString(`clientId`), false)
	if !checkAud {
		go_error.ThrowInternal(`Invalid audience.`)
	}

	// 校验sub
	sub, ok := jwtBody[`sub`]
	if !ok || sub.(string) == `` {
		go_error.ThrowInternal(`Invalid subject.`)
	}

	out.UserId = go_reflect.Reflect.MustToUint64(sub)
	// 校验scope
	jwtScopes, ok := jwtBody[`scope`]
	if !ok {
		go_error.ThrowInternal(`Invalid scope.`)
	}
	if param != nil {
		newParam := param.(OauthJwtValidateParam)
		for _, scope := range newParam.RequiredScopes {
			if !strings.Contains(jwtScopes.(string), scope) {
				go_error.ThrowInternal(`required scope: ` + scope)
			}
		}
	}
}

func (this *OauthJwtValidateClass) getPublicKey(token *jwt.Token) (string, error) {
	cert := ""
	wellKnownParams := &public.WellKnownParams{}
	data, err := global.HydraClientInstance.Public.WellKnown(wellKnownParams.WithTimeout(30 * time.Second))
	if err != nil {
		return ``, err
	}

	for _, v := range data.Payload.Keys {
		if token.Header["kid"].(string) == *v.Kid {
			cert, err = this.getPublickKeyFromJwk(v)
			if err != nil {
				return ``, err
			}
			break
		}
	}
	return cert, nil
}

func (this *OauthJwtValidateClass) getPublickKeyFromJwk(jwkObject *models.JSONWebKey) (string, error) {
	if *jwkObject.Kty != "RSA" {
		return ``, errors.New(`invalid key type`)
	}

	n, err := base64.RawURLEncoding.DecodeString(jwkObject.N)
	if err != nil {
		return ``, err
	}

	e, err := base64.RawURLEncoding.DecodeString(jwkObject.E)
	if err != nil {
		return ``, err
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(n),
		E: int(new(big.Int).SetBytes(e).Int64()),
	}

	der, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		log.Fatal(err)
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}

	var out bytes.Buffer
	err = pem.Encode(&out, block)
	if err != nil {
		return ``, err
	}
	return out.String(), nil
}
