/*
 * Copyright (c) 2024 OceanBase.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-resty/resty/v2"
)

type Response interface {
	GetStatusCode() int
	GetData() interface{}
	GetError() error
	isExpectReturn() bool
}

const (
	DecryptError      int = 1     // 解码错误
	IncompatibleError int = 2     // 校验防护不兼容
	UnauthorizedError int = 10008 // 校验失败
)

// the response struct of ocsagent
type OcsAgentResponse struct {
	Successful bool        `json:"successful"`      // Whether request successful or not
	Timestamp  time.Time   `json:"timestamp"`       // Request handling timestamp (server time)
	Duration   int64       `json:"duration"`        // Request handling time cost (ms)
	Status     int         `json:"status"`          // HTTP status code
	TraceId    string      `json:"traceId"`         // Request trace ID, contained in server logs
	Data       interface{} `json:"data,omitempty"`  // Data payload when response is successful
	Error      *ApiError   `json:"error,omitempty"` // Error payload when response is failed
	ret        bool        // Whether the response is expected to return
}

func (r *OcsAgentResponse) GetData() interface{} {
	return r.Data
}

func (r *OcsAgentResponse) GetError() error {
	return r.Error
}

func (r *OcsAgentResponse) GetStatusCode() int {
	return r.Status
}

func (r *OcsAgentResponse) isExpectReturn() bool {
	return r.ret
}

// the api error struct of ocsagent
type ApiError struct {
	Code      int           `json:"code"`                // Error code
	Message   string        `json:"message"`             // Error message
	SubErrors []interface{} `json:"subErrors,omitempty"` // Sub errors
}

func (a ApiError) Error() string {
	return a.String()
}

func (a ApiError) String() string {
	if len(a.SubErrors) == 0 {
		return fmt.Sprintf("{Code:%v, Message:%v}", a.Code, a.Message)
	} else {
		return fmt.Sprintf("{Code:%v, Message:%v, SubErrors:%+v}", a.Code, a.Message, a.SubErrors)
	}
}

func (a *ApiError) IsError(errorCode int) bool {
	return a.Code == errorCode
}

func Unmarshal(response Response, httpResponse *resty.Response) error {
	if httpResponse.IsError() {
		if response != nil {
			return response.GetError()
		}
		return errors.New("http response error")
	}
	if response.isExpectReturn() {
		responseMap, ok := response.GetData().(map[string]interface{})
		if !ok {
			return errors.New("response data is not map")
		}
		if len(responseMap) == 0 {
			return errors.New("response data is empty")
		}

		data, err := json.Marshal(responseMap)
		if err != nil {
			return err
		}

		if err = unmarshal(data, response); err != nil {
			return err
		}
	}
	return nil
}

func unmarshal(data []byte, ret interface{}) error {
	if reflect.TypeOf(ret).Elem().Kind() == reflect.Slice {
		iterableData := IterableData{}
		err := json.Unmarshal(data, &iterableData)
		if err != nil {
			return err
		}
		data, err = json.Marshal(iterableData.Contents)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(data, ret)
}

// iterable data struct
type IterableData struct {
	Contents interface{} `json:"contents"`
}
