package open

import (
  "github.com/kelvinji2009/antsdk/api"
)

type AlipayOpenPublicShortlinkCreateResponse struct {
  api.AlipayResponse
  ShortLink string `json:"shortlink"` // 生成的带参推广短链接
}
