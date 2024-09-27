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
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type UploadPkgRequest struct {
	*request.BaseRequest
	pkg string
}

// NewUploadPkgRequest return a UploadPkgRequest, which can be used as the argument for the UploadPkgWithRequest.
// params: the parameters to be restored.
func (c *Client) NewUploadPkgRequest(pkg string) *UploadPkgRequest {
	req := &UploadPkgRequest{
		BaseRequest: &request.BaseRequest{},
	}
	req.InitApiInfo("/api/v1/upgrade/package", c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	req.pkg = pkg
	return req
}

func (req *UploadPkgRequest) parsePkg(pkg string) error {
	file := filepath.Base(pkg)
	// 创建缓冲区用以保存multipart的内容
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// 写入文件
	fileWriter, err := bodyWriter.CreateFormFile("file", file)
	if err != nil {
		return err
	}

	fh, err := os.Open(pkg) // 路径替换为你的文件路径
	if err != nil {
		return err
	}
	defer fh.Close()

	// iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	if err := bodyWriter.Close(); err != nil {
		return err
	}
	req.SetBody(bodyBuf.Bytes())
	req.SetHeader("Content-Type", bodyWriter.FormDataContentType())
	return nil
}

type GetUploadPkgResponse struct {
	*response.OcsAgentResponse
	*UpgradePkgInfo
}

type UpgradePkgInfo struct {
	PkgId               int
	Name                string
	Version             string
	ReleaseDistribution string
	Distribution        string
	Release             string
	Architecture        string
	Size                uint64
	PayloadSize         uint64
	ChunkCount          int
	Md5                 string
	UpgradeDepYaml      string
	GmtModify           time.Time
}

func (c *Client) createUploadPkgResponse() *GetUploadPkgResponse {
	resp := &GetUploadPkgResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		UpgradePkgInfo:   &UpgradePkgInfo{},
	}
	resp.Data = resp.UpgradePkgInfo
	return resp
}

// UploadPkg returns a error, when upload package successfully, the error will be nil.
// params: the parameters to be restored.
func (c *Client) UploadPkg(pkg string) (UploadPkgResp *UpgradePkgInfo, err error) {
	req := c.NewUploadPkgRequest(pkg)
	return c.UploadPkgWithRequest(req)
}

// UploadPkgWithRequest returns an error, when upload package successfully, the error will be nil.
// the parameter is a UploadPkgRequest, which can be created by NewUploadPkgRequest.
func (c *Client) UploadPkgWithRequest(req *UploadPkgRequest) (UploadPkgResp *UpgradePkgInfo, err error) {
	if err := req.parsePkg(req.pkg); err != nil {
		return nil, err
	}
	response := c.createUploadPkgResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.UpgradePkgInfo, err
}
