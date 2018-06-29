package api

import (
  "github.com/kelvinji2009/antsdk/utils"
)

type IAlipayUploadRequest interface {
  IAlipayRequest
  GetFileParams() map[string]*utils.FileItem
}
