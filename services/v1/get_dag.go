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

package v1

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/obshell-sdk-go/model"
	"github.com/obshell-sdk-go/sdk/request"
	"github.com/obshell-sdk-go/sdk/response"
)

var ErrQueryDagFailed = errors.New("query dag failed")

type GetDagRequest struct {
	*request.BaseRequest
}

// NewGetDagRequest return a GetDagRequest, which can be used as the argument for the GetDagWithRequest.
// You can set whether show detail by calling SetShowDetail.
func (c *Client) NewGetDagRequest(dagId string) *GetDagRequest {
	req := &GetDagRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/task/dag/"+dagId, c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

// SetShowDetail set whether show detail.
func (r *GetDagRequest) SetShowDetail(showDetail bool) *GetDagRequest {
	r.SetBody(map[string]bool{"showDetail": showDetail})
	return r
}

type GetDagResponse struct {
	*response.TaskResponse
}

func (c *Client) createGetDagResponse() *GetDagResponse {
	return &GetDagResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// GetDag returns a DagDetailDTO and an error.
// If the error is non-nil, the DagDetailDTO will be nil.
// dagId is the id of the dag.
// If you don't want to show detail, you need to use NewGetDagRequest and call SetShowDetail(false).
func (c *Client) GetDag(dagId string) (*model.DagDetailDTO, error) {
	req := c.NewGetDagRequest(dagId)
	return c.GetDagWithRequest(req)
}

// GetDagWithRequest returns a DagDetailDTO and an error.
// The parameter is a GetDagRequest, which can be created by NewGetDagRequest.
// If the error is non-nil, the DagDetailDTO will be nil.
func (c *Client) GetDagWithRequest(req *GetDagRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createGetDagResponse()
	err = c.Execute(req, response)
	return response.DagDetailDTO, err
}

// WaitDagSucceed wait for a dag to succeed(return error if the dag is failed or occur error when query dag).
// When query dag failed, the error will be wrapped with v1.ErrQueryDagFailed.
// Return err once a query failed
func (c *Client) WaitDagSucceed(dagId string) (dag *model.DagDetailDTO, err error) {
	for {
		dag, err = c.GetDag(dagId)
		if err != nil {
			return nil, errors.Wrap(ErrQueryDagFailed, err.Error())
		}
		if dag.IsSucceed() {
			return
		}
		if dag.IsFailed() {
			logs := model.GetFailedDagLastLog(dag)
			err = errors.New(strings.Join(logs, "\n"))
			return
		}
		time.Sleep(2 * time.Second)
	}
}

// WaitDagSucceed wait for a dag to succeed(return error if the dag is failed or occur error when query dag).
// When query dag failed, WaitDagSucceedWithRetry will retry until the dag is finished or the retry times has reached the limit.
// When query dag failed, the error will be wrapped with v1.ErrQueryDagFailed.
func (c *Client) WaitDagSucceedWithRetry(dagId string, retryTimes int) (dag *model.DagDetailDTO, err error) {
	for {
		dag, err = c.GetDag(dagId)
		if err != nil {
			if retryTimes != 0 {
				retryTimes--
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, errors.Wrap(ErrQueryDagFailed, err.Error())
		}
		if dag.IsSucceed() {
			return
		}
		if dag.IsFailed() {
			logs := model.GetFailedDagLastLog(dag)
			err = errors.New(strings.Join(logs, "\n"))
			return
		}
		time.Sleep(2 * time.Second)
	}
}
