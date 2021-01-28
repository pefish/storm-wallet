package global

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/ory/hydra-client-go/client"
	go_config "github.com/pefish/go-config"
	go_http "github.com/pefish/go-http"
	go_logger "github.com/pefish/go-logger"
	"github.com/pkg/errors"
)

var (
	HydraClientInstance *client.OryHydra
	AuthServerUrl       string
)

type global struct {
	Seeds map[string]string
}

var Global = global{}

func (g *global) Init() {
	seeds, err := g.getAllSeeds(go_config.Config.MustGetString(`privateManagerUrl`))
	if err != nil {
		panic(err)
	}
	go_logger.Logger.Debug(`all seeds is loaded`)
	Global.Seeds = seeds
}

type ApiResult struct {
	Code uint8
	Msg  string
	Data map[string]map[string]string
}

func (g *global) getAllSeeds(url string) (map[string]string, error) {
	var result ApiResult
	_, err := go_http.NewHttpRequester().
		PostForStruct(go_http.RequestParam{
			Url: url + `/seed/get_all_seeds`,
		}, &result)
	if err != nil {
		return nil, errors.Wrap(err, `get all seed error`)
	}
	return result.Data[`seeds`], nil
}

func AesCbcDecrypt(key string, data string) (string, error) {
	length := len(key)
	if length <= 16 {
		key = spanLeft(key, 16, `0`)
	} else if length <= 24 {
		key = spanLeft(key, 24, `0`)
	} else if length <= 32 {
		key = spanLeft(key, 32, `0`)
	} else {
		return "", errors.New(`length of secret key error`)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, make([]byte, blockSize))
	crypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origDataLength := len(origData)
	unpadding := int(origData[origDataLength-1])
	origData = origData[:(origDataLength - unpadding)]
	return string(origData), nil
}

func spanLeft(str string, length int, fillChar string) string {
	result := ``
	for i := 0; i < length-len(str); i++ {
		result += fillChar
	}
	return result + str
}
