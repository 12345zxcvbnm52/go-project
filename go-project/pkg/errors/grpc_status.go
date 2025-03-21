package errors

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	grpccode "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	UseGobMarshaler  = 1
	UseJsonMarshaler = 2
)

var marshalerCtn = UseJsonMarshaler

func SetErrorMarshalType(mte int) {
	switch mte {
	case UseGobMarshaler:
		marshalerCtn = UseGobMarshaler
	case UseJsonMarshaler:
		marshalerCtn = UseJsonMarshaler
	default:
		marshalerCtn = UseJsonMarshaler
	}
}

type codePayload struct {
	Code     int    `json:"code"`
	HttpCode int    `json:"http_code"`
	GrpcCode int    `json:"grpc_code"`
	Message  string `json:"message"`
}

type marshalData struct {
	CodeMsg  codePayload `json:"code_msg"`
	StackMsg string      `json:"stack_msg"`
}

func (c *withCode) marshalJSON() ([]byte, error) {
	data := marshalData{
		CodeMsg: codePayload{
			Code:     c.code.ErrorCode(),
			HttpCode: c.code.HTTPCode(),
			GrpcCode: int(c.code.GrpcCode()),
			Message:  c.Message(),
		},
		StackMsg: fmt.Sprintf("%+v", c),
	}

	return json.Marshal(data)
}

func (c *withCode) marshalGOB() (data []byte, err error) {
	var buf bytes.Buffer
	msr := gob.NewEncoder(&buf)

	value := marshalData{
		CodeMsg: codePayload{
			Code:     c.code.ErrorCode(),
			HttpCode: c.code.HTTPCode(),
			GrpcCode: int(c.code.GrpcCode()),
			Message:  c.Message(),
		},
		StackMsg: fmt.Sprintf("%+v", c),
	}
	err = msr.Encode(value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *withCode) unmarshalJSON(data []byte) error {
	var temp marshalData
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	c.code = &defaultCoder{
		code:     temp.CodeMsg.Code,
		httpCode: temp.CodeMsg.HttpCode,
		grpcCode: grpccode.Code(temp.CodeMsg.GrpcCode),
		message:  temp.CodeMsg.Message,
	}

	c.cause = &fundamental{
		msg:   temp.StackMsg,
		stack: callers(),
	}

	return nil
}

func (c *withCode) unmarshalGOB(data []byte) error {
	var temp marshalData
	evt := gob.NewDecoder(bytes.NewBuffer(data))
	if err := evt.Decode(&temp); err != nil {
		return err
	}

	c.code = &defaultCoder{
		code:     temp.CodeMsg.Code,
		httpCode: temp.CodeMsg.HttpCode,
		grpcCode: grpccode.Code(temp.CodeMsg.GrpcCode),
		message:  temp.CodeMsg.Message,
	}

	c.cause = &fundamental{
		msg:   temp.StackMsg,
		stack: callers(),
	}
	return nil
}

func (e *withCode) grpcStatus() *status.Status {
	var msg []byte
	var err error
	switch marshalerCtn {
	case UseGobMarshaler:
		msg, err = e.marshalGOB()
	case UseJsonMarshaler:
		msg, err = e.marshalJSON()
	}

	if err != nil {
		// 若序列化失败，返回原始错误信息，避免 status.New 参数为空
		return status.New(e.code.GrpcCode(), fmt.Sprintf("failed to marshal error: %v", err))
	}
	return status.New(e.code.GrpcCode(), string(msg))
}

func UnmarshalCodeError(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	e := &withCode{stack: callers()}
	switch marshalerCtn {
	case UseGobMarshaler:
		e.unmarshalGOB(data)
	case UseJsonMarshaler:
		e.unmarshalJSON(data)
	}
	return e
}

func MarshalCodeError(err error) string {
	if cerr, ok := err.(*withCode); ok {
		data, _ := cerr.marshalJSON()
		return string(data)
	} else {
		err = WithCoder(err, CodeInternalError, "")
		data, _ := err.(*withCode).marshalJSON()
		return string(data)
	}
}

// 从gRPC错误提取withCode结构
func ExtractCodeErrorFromGRPC(err error) error {
	if st, ok := status.FromError(err); ok {
		var c withCode
		if jsonErr := json.Unmarshal([]byte(st.Message()), &c); jsonErr != nil {
			return err
		}
		return &c
	}
	return err
}

func (w *withCode) GRPCStatus() *status.Status {
	return w.grpcStatus()
}
