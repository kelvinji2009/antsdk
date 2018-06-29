package open

import (
  "github.com/kelvinji2009/antsdk/api"
)

type AlipayOpenPublicMessageTotalSendResponse struct {
  api.AlipayResponse
  MessageId string `json:"message_id"`  // 消息ID
}
