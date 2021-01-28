package address

import (
	"fmt"
	go_config "github.com/pefish/go-config"
	go_logger "github.com/pefish/go-logger"
	"github.com/pkg/errors"
	"strconv"
	"wallet-storm-wallet/global"
	"wallet-storm-wallet/util/kto"
)

type DepositAddressClass struct {
	baseUrl   string
	apiConfig map[string]interface{}
}

var DepositAddressService = DepositAddressClass{}

func (dac *DepositAddressClass) ValidateAddress(series string, address string, memo string) (bool, error) {
	return true, nil
}

type GetAddressReturn struct {
	Address string `json:"address"`
	Path    string `json:"path"`
}

func getCurrencySeed(series string) (string, error) {
	encryptedSeed, ok := global.Global.Seeds[series]
	if !ok {
		return "", errors.New(series + `seed not exist or config error`)
	}
	return global.AesCbcDecrypt(go_config.Config.MustGetString(`secret`), encryptedSeed)
}

func (dac *DepositAddressClass) GetAddress(series string, type_ uint64, index uint64) (*GetAddressReturn, error) {
	switch series {
	case `Kto`:
		seed, err := getCurrencySeed(series)
		if err != nil {
			return nil, err
		}
		path := `m/` + strconv.Itoa(int(type_)) + `/` + strconv.FormatUint(index, 10)
		go_logger.Logger.DebugF(`series: %s path: %s`, series, path)
		address, _, err := kto.GenKeyPairWithSeedAndPath(seed, path)
		if err != nil {
			if err.Error() == `address size error` {
				path += `/` + strconv.FormatUint(index, 10)
				address, _, err = kto.GenKeyPairWithSeedAndPath(seed, path)
			}
		}
		return &GetAddressReturn{
			Address: address,
			Path:    path,
		}, nil
	default:
		return nil, fmt.Errorf("series: %s not supported", series)
	}
}
