package zhima

import (
  "github.com/kelvinji2009/antsdk/api"
)

type ZhimaDataFeedbackurlQueryResponse struct {
  api.AlipayResponse
  FeedbackUrl string `json:"feedback_url"`  // 反馈模板地址
}
