package external_service

import (
	"github.com/pefish/go-config"
	external_service "github.com/pefish/go-core/driver/external-service"
)

type DepositAddressClass struct {
	baseUrl   string
	apiConfig map[string]interface{}
	external_service.BaseExternalServiceClass
}

var DepositAddressService = DepositAddressClass{}

func (this *DepositAddressClass) Init(driver *external_service.ExternalServiceDriverClass) {
	this.apiConfig = go_config.Config.MustGetMap(`depositAddressApi`)
	this.baseUrl = this.apiConfig[`baseUrl`].(string)
}

func (this *DepositAddressClass) ValidateAddress(series string, address string, memo string) {
	path := this.apiConfig[`validateAddressPath`].(string)
	this.PostJson(this.baseUrl+path, map[string]interface{}{
		`series`:  series,
		`address`: address,
		`memo`:    memo,
	})
}

type GetAddressReturn struct {
	Address string `json:"address"`
	Path    string `json:"path"`
}

func (this *DepositAddressClass) GetAddress(series string, type_ uint64, index uint64) GetAddressReturn {
	path := this.apiConfig[`getAddressPath`].(string)
	result := GetAddressReturn{}
	this.PostJsonForStruct(this.baseUrl+path, map[string]interface{}{
		`series`: series,
		`type`:   type_,
		`index`:  index,
	}, &result)
	return result
}
