package comm

import "context"

func GetReqId(ctx context.Context) (reqId string) {
	if requestID, ok := ctx.Value("reqId").(string); ok {
		return requestID
	}
	return "EMPTY_REQUEST_ID"
}
