package data

import (
  "github.com/kelvinji2009/antsdk/api"
)

type KoubeiMarketingDataCustomreportDetailQueryResponse struct {
  api.AlipayResponse
  ReportConditionInfo CustomReportCondition `json:"report_condition_info"`  // 自定义报表规则条件的详细信息
}
