//MIT License
//
//Copyright (c) 2016 GenesisKernel
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.
package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/GenesisCommunity/go-genesis/packages/converter"
	"github.com/GenesisCommunity/go-genesis/packages/script"
	"github.com/GenesisCommunity/go-genesis/packages/utils/tx"

	log "github.com/sirupsen/logrus"
)

type prepareResult struct {
	ForSign string            `json:"forsign"`
	Signs   []TxSignJSON      `json:"signs"`
	Values  map[string]string `json:"values"`
	Time    string            `json:"time"`
}

func prepareContract(w http.ResponseWriter, r *http.Request, data *apiData, logger *log.Entry) error {
	var (
		result  prepareResult
		timeNow int64
		smartTx tx.SmartContract
	)

	timeNow = time.Now().Unix()
	result.Time = converter.Int64ToStr(timeNow)
	result.Values = make(map[string]string)
	contract, parerr, err := validateSmartContract(r, data, &result)
	if err != nil {
		if strings.HasPrefix(err.Error(), `E_`) {
			return errorAPI(w, err.Error(), http.StatusBadRequest, parerr)
		}
		return errorAPI(w, err, http.StatusBadRequest)
	}
	info := (*contract).Block.Info.(*script.ContractInfo)
	smartTx.TokenEcosystem = data.params[`token_ecosystem`].(int64)
	smartTx.MaxSum = data.params[`max_sum`].(string)
	smartTx.PayOver = data.params[`payover`].(string)
	if data.params[`signed_by`] != nil {
		smartTx.SignedBy = data.params[`signed_by`].(int64)
	}
	smartTx.Header = tx.Header{Type: int(info.ID), Time: timeNow, EcosystemID: data.ecosystemId, KeyID: data.keyId}
	forsign := smartTx.ForSign()
	if info.Tx != nil {
		for _, fitem := range *info.Tx {
			if strings.Contains(fitem.Tags, `image`) || strings.Contains(fitem.Tags, `signature`) {
				continue
			}
			var val string
			if fitem.Type.String() == `[]interface {}` {
				for key, values := range r.Form {
					if key == fitem.Name+`[]` {
						var list []string
						for _, value := range values {
							list = append(list, value)
						}
						val = strings.Join(list, `,`)
					}
				}
			} else {
				val = strings.TrimSpace(r.FormValue(fitem.Name))
				if strings.Contains(fitem.Tags, `address`) {
					val = converter.Int64ToStr(converter.StringToAddress(val))
				} else if fitem.Type.String() == script.Decimal {
					val = strings.TrimLeft(val, `0`)
				} else if fitem.Type.String() == `int64` && len(val) == 0 {
					val = `0`
				}
			}
			forsign += fmt.Sprintf(",%v", val)
		}
	}
	result.ForSign = forsign
	data.result = result
	return nil
}
