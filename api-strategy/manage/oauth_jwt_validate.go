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
	_type "github.com/pefish/go-core/api-session/type"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-reflect"
	"log"
	"math/big"
	"strings"
	"time"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/global"

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

func (this *OauthJwtValidateClass) InitAsync(param interface{}, onAppTerminated chan interface{}) {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s InitAsync`, this.GetName())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s InitAsync defer`, this.GetName())
}

func (this *OauthJwtValidateClass) Init(param interface{}) {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init`, this.GetName())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init defer`, this.GetName())
}

type OauthJwtValidateParam struct {
	RequiredScopes []string
}

func (this *OauthJwtValidateClass) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	// 校验jwt合法性
	jwtStr := out.Header(`JSON-WEB-TOKEN`)
	if jwtStr == `` {
		return go_error.Wrap(errors.New(`auth error. jwt not found.`))
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
		return go_error.Wrap(errors.New(`jwt is illegal`))
	}

	jwtBody := token.Claims.(jwt.MapClaims)
	out.SetJwtBody(jwtBody)
	// 校验iss
	checkIss := jwtBody.VerifyIssuer(global.AuthServerUrl+`/`, false)
	if !checkIss {
		return go_error.Wrap(errors.New(`Invalid issuer.`))
	}

	// 校验aud[0]必须是clientId
	checkAud := jwtBody.VerifyAudience(go_config.Config.MustGetString(`clientId`), false)
	if !checkAud {
		return go_error.Wrap(errors.New(`Invalid audience.`))
	}

	// 校验sub
	sub, ok := jwtBody[`sub`]
	if !ok || sub.(string) == `` {
		return go_error.Wrap(errors.New(`Invalid subject.`))
	}
	out.SetUserId(go_reflect.Reflect.MustToUint64(sub))
	// 校验scope
	jwtScopes, ok := jwtBody[`scope`]
	if !ok {
		return go_error.Wrap(errors.New(`Invalid scope.`))
	}
	if param != nil {
		newParam := param.(OauthJwtValidateParam)
		for _, scope := range newParam.RequiredScopes {
			if !strings.Contains(jwtScopes.(string), scope) {
				return go_error.Wrap(errors.New(`required scope: ` + scope))
			}
		}
	}
	return nil
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
