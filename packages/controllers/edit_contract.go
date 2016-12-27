// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package controllers

import (
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
)

type editContractPage struct {
	Alert        string
	SignData     string
	ShowSignData bool
	CountSignArr []int
	Lang         map[string]string
	Data         map[string]string
	WalletId     int64
	CitizenId    int64
	TxType       string
	TxTypeId     int64
	TimeNow      int64
	TableName    string
	StateId      int64
	Global       string
}

func (c *Controller) EditContract() (string, error) {

	txType := "EditContract"
	txTypeId := utils.TypeInt(txType)
	timeNow := utils.Time()

	global := c.r.FormValue("global")
	prefix := "global"
	if global == "" || global == "0" {
		prefix = c.StateIdStr
		global = "0"
	}

	id := utils.StrToInt64(c.r.FormValue("id"))
	name := c.r.FormValue("name")
	if len(name) > 0 && name[:1] == `@` {
		name = name[1:]
	}
	if len(name) > 0 && !utils.CheckInputData_(name, "string", "") {
		return "", utils.ErrInfo("Incorrect name")
	}
	var data map[string]string
	var err error
	if id != 0 {
		data, err = c.OneRow(`SELECT * FROM "`+prefix+`_smart_contracts" WHERE id = ?`, id).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	} else {
		data, err = c.OneRow(`SELECT * FROM "`+prefix+`_smart_contracts" WHERE name = ?`, name).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	TemplateStr, err := makeTemplate("edit_contract", "editContract", &editContractPage{
		Alert:        c.Alert,
		Lang:         c.Lang,
		ShowSignData: c.ShowSignData,
		SignData:     "",
		WalletId:     c.SessWalletId,
		Data:         data,
		Global:       global,
		CitizenId:    c.SessCitizenId,
		CountSignArr: c.CountSignArr,
		TimeNow:      timeNow,
		TxType:       txType,
		TxTypeId:     txTypeId,
		StateId:      c.SessStateId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
