package xxl

import (
	"encoding/json"
	"net"
	"strconv"
)

func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

func LocalIP() string {
	ips, _ := LocalIPv4s()
	return ips[0]
}

// Int64ToStr int64 to str
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

//执行任务回调
func returnCall(req *RunReq, code int64, msg string) []byte {
	data := call{
		&callElement{
			LogID:      req.LogID,
			LogDateTim: req.LogDateTime,
			ExecuteResult: &ExecuteResult{
				Code: code,
				Msg:  msg,
			},
			HandleCode: int(code),
			HandleMsg:  msg,
		},
	}
	str, _ := json.Marshal(data)
	return str
}

//杀死任务返回
func returnKill(req *killReq, code int64) []byte {
	msg := ""
	if code != SuccessCode {
		msg = "Task kill err"
	}
	data := res{
		Code: code,
		Msg:  msg,
	}
	str, _ := json.Marshal(data)
	return str
}

//忙碌返回
func returnIdleBeat(code int64) []byte {
	msg := ""
	if code != SuccessCode {
		msg = "Task is busy"
	}
	data := res{
		Code: code,
		Msg:  msg,
	}
	str, _ := json.Marshal(data)
	return str
}

//通用返回
func returnGeneral() []byte {
	data := &res{
		Code: SuccessCode,
		Msg:  "",
	}
	str, _ := json.Marshal(data)
	return str
}
