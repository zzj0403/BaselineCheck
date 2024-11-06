package baselinelinux

import (
	"BaselineCheck/client/comm"
	"BaselineCheck/client/getinfo"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Result struct {
	BaseInfo getinfo.BaseInfo `json:"base_info"`
	// TrojanInfo     []getinfo.SingleTrojanInfo     `json:"trojan_info"`
	// EmergencyInfo  getinfo.EmergencyInfo          `json:"emergency_info"`
	ComplianceInfo []getinfo.SingleComplianceInfo `json:"compliance_info"`
}

func Run(pushUrl string) {
	var r Result
	var bs getinfo.BaseInfo
	bs.GetBaseInfo()
	r.BaseInfo = bs

	// r.TrojanInfo = getinfo.ReturnResultTro()

	// var eg getinfo.EmergencyInfo
	// eg.GetEmergencyInfo()
	// r.EmergencyInfo = eg

	r.ComplianceInfo = getinfo.ReturnResultCom()

	r.BaseInfo.Description = "基线检查任务"
	r.BaseInfo.End = int(time.Now().Unix())

	res, err := json.MarshalIndent(r, "", "  ") // 格式化编码
	if err != nil {
		log.Println("JSON ERR:", err)
	}
	log.Println("[✓] Baseline check finish!")
	comm.JsonWrite(res) // 把结果写入文件
	// 把结果发送到 api 服务器
	resp, err := http.Post(pushUrl, "application/json", bytes.NewReader(res))
	if err != nil {
		log.Println("POST ERR:", err)
	}
	// 打印返回结果
	defer resp.Body.Close()

	// 读取并打印响应
	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Println("Read response body ERR:", err)
		return
	}

	// 打印响应内容
	log.Printf("Response Status: %s\n", resp.Status)
	log.Printf("Response Body: %s\n", body.String())
}

func (r *Result) LengthComplianceInfo() int {
	// 统计ComplianceInfo列表里面有多少个
	return len(r.ComplianceInfo)
}
