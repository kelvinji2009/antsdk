package open

import (
  "github.com/kelvinji2009/antsdk/api"
)

type AlipayOpenPublicMenuQueryResponse struct {
  api.AlipayResponse
  MenuContent string `json:"menu_content"`  // 一级菜单数组，个数应为1~4个
}
