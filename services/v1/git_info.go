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
	"github.com/obshell-sdk-go/sdk/request"
	"github.com/obshell-sdk-go/sdk/response"
)

type GetGitInfoRequest struct {
	*request.BaseRequest
}

type GitInfo struct {
	GitBranch        string `json:"branch"`
	GitCommitId      string `json:"commitId"`
	GitShortCommitId string `json:"shortCommitId"`
	GitCommitTime    string `json:"commitTime"`
}

// NewGetGitInfoRequest return a GetGitInfoRequest, which can be used as the argument for the GetGitInfoWithRequest.
func (c *Client) NewGetGitInfoRequest() *GetGitInfoRequest {
	req := &GetGitInfoRequest{
		BaseRequest: &request.BaseRequest{},
	}
	req.InitApiInfo("/api/v1/git-info", c.GetHost(), c.GetPort(), "GET")
	return req
}

type GetGitInfoResponse struct {
	*response.OcsAgentResponse
	*GitInfo
}

func (c *Client) createGetGitInfoResponse() *GetGitInfoResponse {
	return &GetGitInfoResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
}

// GetGitInfo returns a ObGitInfoResp and an error.
// If the error is non-nil, the ObGitInfoResp will be nil.
func (c *Client) GetGitInfo() (ObGitInfoResp *GitInfo, err error) {
	req := c.NewGetGitInfoRequest()
	return c.GetGitInfoWithRequest(req)
}

// GetGitInfoWithRequest returns a ObGitInfoResp and an error.
// If the error is non-nil, the ObGitInfoResp will be nil.
func (c *Client) GetGitInfoWithRequest(req *GetGitInfoRequest) (ObGitInfoResp *GitInfo, err error) {
	response := c.createGetGitInfoResponse()
	err = c.Execute(req, response)
	return response.GitInfo, err
}
