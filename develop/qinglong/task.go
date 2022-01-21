package qinglong

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core"
)

func initTask() {
	core.AddCommand("ql", []core.Function{
		{
			Rules: []string{`\r\a\w ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				s.UAtLast()
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron := &Carrier{
						Get: "data._id",
					}
					_, err := Req(ql, cron, CRONS, POST, []byte(`{"name":"sillyGirl临时创建任务","command":"ql raw `+s.Get()+`","schedule":" 1 1 1 1 1"}`))
					if err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					s.Reply(fmt.Sprintf("任务ID%s。", cron.Value))
					if _, err := Req(ql, CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.Value))); err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					if err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					for {
						time.Sleep(time.Microsecond * 300)
						data, _ := GetCronLog(ql, cron.Value)
						if strings.Contains(data, "执行结束...") {
							s.Reply(data + ql.GetTail())
							break
						}
					}
					if _, err := Req(ql, CRONS, DELETE, []byte(`["`+cron.Value+`"]`)); err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
				}
				return nil
			},
		},
		{
			Rules: []string{`task ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron := &Carrier{
						Get: "data._id",
					}
					_, err := Req(ql, cron, CRONS, POST, []byte(`{"name":"sillyGirl临时创建任务","command":"task `+s.Get()+`","schedule":" 1 1 1 1 1"}`))
					if err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					if _, err := Req(ql, CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.Value))); err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					for {
						time.Sleep(time.Second)
						data, _ := GetCronLog(ql, cron.Value)
						if strings.Contains(data, "执行结束...") {
							s.Reply(data + ql.GetTail())
							break
						}
					}
					if _, err := Req(ql, cron, CRONS, DELETE, []byte(`["`+cron.Value+`"]`)); err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
				}
				return nil
			},
		},
		{
			Rules: []string{`repo ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				err, qls := QinglongSC(s)
				if err != nil {
					return err
				}
				for _, ql := range qls {
					cron := &Carrier{
						Get: "data._id",
					}
					data, _ := json.Marshal(map[string]string{
						"name":     "sillyGirl临时创建任务",
						"command":  `ql repo ` + s.Get(),
						"schedule": "1 1 1 1 1",
					})
					_, err := Req(ql, cron, CRONS, POST, data)
					if err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					if _, err := Req(ql, CRONS, PUT, "/run", []byte(fmt.Sprintf(`["%s"]`, cron.Value))); err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
					for {
						time.Sleep(time.Second)
						data, _ := GetCronLog(ql, cron.Value)
						if strings.Contains(data, "执行结束...") {
							s.Reply(data + ql.GetTail())
							break
						}
					}
					if _, err := Req(ql, cron, CRONS, DELETE, []byte(`["`+cron.Value+`"]`)); err != nil {
						s.Reply(err.Error() + ql.GetTail())
						continue
					}
				}
				return nil
			},
		},
	})
}
