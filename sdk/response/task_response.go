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

import "github.com/oceanbase/obshell-sdk-go/model"

type TaskResponse struct {
	*OcsAgentResponse
	*model.DagDetailDTO
}

func NewTaskResponse() *TaskResponse {
	resp := &TaskResponse{
		OcsAgentResponse: NewOcsAgentResponse(),
		DagDetailDTO:     &model.DagDetailDTO{},
	}
	resp.Data = resp.DagDetailDTO
	return resp
}

func NewOcsAgentResponse() *OcsAgentResponse {
	resp := &OcsAgentResponse{
		ret: true,
	}
	return resp
}

func NewOcsAgentResponseWithoutReturn() *OcsAgentResponse {
	resp := &OcsAgentResponse{}
	return resp
}

type IterableDagDetailDTO struct {
	Contents []*model.DagDetailDTO `json:"contents"`
}

type MutilTaskReponse struct {
	*OcsAgentResponse
	*IterableDagDetailDTO
}

func NewMutilTaskReponse() *MutilTaskReponse {
	resp := &MutilTaskReponse{
		OcsAgentResponse:     NewOcsAgentResponse(),
		IterableDagDetailDTO: &IterableDagDetailDTO{},
	}
	resp.Data = resp.IterableDagDetailDTO
	return resp
}
