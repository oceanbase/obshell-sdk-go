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
	"fmt"

	"github.com/pkg/errors"

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type OperateDagRequest struct {
	*request.BaseRequest
	id       string
	operator string
}

// NewOperateDagRequest return a OperateDagRequest, which can be used as the argument for the OperateDagWithRequest/OperateDagSyncWithRequest.
// dagId: the id of the dag to be operated.
// operator: the operator of the dag to be operated, it can be "PASS", "ROLLBACK", "RETRY" or "CANCEL".
func (c *Client) NewOperateDagRequest(dagId string, operator string) *OperateDagRequest {
	req := &OperateDagRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		id:          dagId,
		operator:    operator,
	}
	req.InitApiInfo("/api/v1/task/dag/"+dagId, c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	req.SetBody(
		map[string]string{
			"operator": operator,
		},
	)
	return req
}

type OperateDagResponse struct {
	*response.OcsAgentResponse
}

func (r *OperateDagResponse) Init() {}

func (c *Client) createOperateDagResponse() *OperateDagResponse {
	return &OperateDagResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// OperateDag returns an error, when operate dag is completed successfully, the error will be nil.
// dagId: the id of the dag to be operated.
// operator: the operator of the dag to be operated, it can be "PASS", "ROLLBACK", "RETRY" or "CANCEL".
func (c *Client) OperateDag(dagId string, operator string) error {
	req := c.NewOperateDagRequest(dagId, operator)
	return c.OperateDagSyncWithRequest(req)
}

// OperateDagWithRequest returns an error, when dag operator is requested successfully, the error will be nil.
// the parameter is a OperateDagRequest, which can be created by NewOperateDagRequest.
// You can use WaitDagSucceed to wait for the task to complete.
func (c *Client) OperateDagWithRequest(request *OperateDagRequest) (err error) {
	response := c.createOperateDagResponse()
	return c.Execute(request, response)
}

// OperateDagSyncWithRequest returns an error, when the dag operator task is completed successfully, the error will be nil.
// the parameter is a OperateDagRequest, which can be created by NewOperateDagRequest.
func (c *Client) OperateDagSyncWithRequest(request *OperateDagRequest) error {
	if err := c.OperateDagWithRequest(request); err != nil {
		return errors.Wrap(err, "Error occured when operating dag")
	}
	if request.operator == model.PASS_STR { // needn't schedule
		return nil
	}
	dag, err := c.GetDag(request.id)
	if err != nil {
		return fmt.Errorf("Error occured when watching dag %s", dag.GenericID)
	}
	dag, err = c.WaitDagSucceed(dag.GenericID)
	switch request.operator {
	case model.ROLLBACK_STR:
		if dag.IsSucceed() && dag.IsRollback() {
			return nil
		}
	case model.RETRY_STR:
		if dag.IsSucceed() && dag.IsRun() {
			return nil
		}
	case model.CANCEL_STR:
		if dag.IsFailed() && dag.IsCancel() {
			return nil
		}
	}
	return err
}
