// Copyright 2019 (c) Microsoft Corporation.
// Licensed under the MIT license.

//
// Author:
//      Auto Generated on 9/18/2020 using wmigen
//      Source root.CIMV2
//////////////////////////////////////////////
package cimv2

import (
	"github.com/microsoft/wmi/pkg/base/query"
	cim "github.com/microsoft/wmi/pkg/wmiinstance"
)

// Win32_UTCTime struct
type Win32_UTCTime struct {
	*Win32_CurrentTime
}

func NewWin32_UTCTimeEx1(instance *cim.WmiInstance) (newInstance *Win32_UTCTime, err error) {
	tmp, err := NewWin32_CurrentTimeEx1(instance)

	if err != nil {
		return
	}
	newInstance = &Win32_UTCTime{
		Win32_CurrentTime: tmp,
	}
	return
}

func NewWin32_UTCTimeEx6(hostName string,
	wmiNamespace string,
	userName string,
	password string,
	domainName string,
	query *query.WmiQuery) (newInstance *Win32_UTCTime, err error) {
	tmp, err := NewWin32_CurrentTimeEx6(hostName, wmiNamespace, userName, password, domainName, query)

	if err != nil {
		return
	}
	newInstance = &Win32_UTCTime{
		Win32_CurrentTime: tmp,
	}
	return
}
