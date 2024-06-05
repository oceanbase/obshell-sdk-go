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
	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type GetNodeRequest struct {
	*request.BaseRequest
}

// NewGetNodeRequest return a GetNodeRequest, which can be used as the argument for the GetNodeWithRequest.
// You can set whether show detail by calling SetShowDetail.
func (c *Client) NewGetNodeRequest(nodeId string) *GetNodeRequest {
	req := &GetNodeRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/task/node/"+nodeId, c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

// SetShowDetail set whether show detail.
func (r *GetNodeRequest) SetShowDetail(showDetail bool) *GetNodeRequest {
	r.SetBody(map[string]bool{"showDetail": showDetail})
	return r
}

type GetNodeResponse struct {
	*response.OcsAgentResponse
	*model.NodeDetailDTO
}

func (c *Client) createGetNodeResponse() *GetNodeResponse {
	return &GetNodeResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
}

// GetNode returns a NodeDetailDTO and an error.
// If the error is non-nil, the NodeDetailDTO will be nil.
// nodeId is the id of the node.
// If you don't want to show detail, you need to use NewGetNodeRequest and call SetShowDetail(false).
func (c *Client) GetNode(nodeId string) (*model.NodeDetailDTO, error) {
	request := c.NewGetNodeRequest(nodeId)
	return c.GetNodeWithRequest(request)
}

// GetNodeWithRequest returns a NodeDetailDTO and an error.
// The parameter is a GetNodeRequest, which can be created by NewGetNodeRequest.
// If the error is non-nil, the NodeDetailDTO will be nil.
func (c *Client) GetNodeWithRequest(req *GetNodeRequest) (node *model.NodeDetailDTO, err error) {
	response := c.createGetNodeResponse()
	err = c.Execute(req, response)
	return response.NodeDetailDTO, err
}
