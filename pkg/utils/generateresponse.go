package utils

import "go.uber.org/zap"

type MessageResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Response any    `json:"result"`
}

type ErrorResponse struct {
	ErrorCode uint   `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func GenerateResponse(message *string, bodyResponse any) any {
	zap.L().Debug("message response", zap.Any("body", bodyResponse))
	var responseDone MessageResponse
	responseDone.Success = true
	if message != nil {
		responseDone.Message = *message
	}
	if bodyResponse != nil {
		responseDone.Response = bodyResponse
	}
	return responseDone
}
