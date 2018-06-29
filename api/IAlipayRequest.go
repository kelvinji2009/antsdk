package api

import (
  "github.com/kelvinji2009/antsdk/utils"
)

type IAlipayRequest interface {
  GetApiMethodName() string
  GetApiVersion() string
  GetTerminalType() string
  GetTerminalInfo() string
  GetNotifyUrl() string
  GetReturnUrl() string
  GetProdCode() string
  GetTextParams() *utils.AlipayHashMap
  IsNeedEncrypt() bool
}
