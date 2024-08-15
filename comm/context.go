package comm

import "context"

func GetReqId(c context.Context) (reqId string) {
	if v := c.Value("req_id"); v != nil {
		reqId = v.(string)
	}
	return
}
