package address

import (
	"fmt"
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

func (dac *DepositAddressClass) GetAddress(series string, type_ uint64, index uint64) (*GetAddressReturn, error) {
	switch series {
	default:
		return nil, fmt.Errorf("series: %s not supported", series)
	}
}
