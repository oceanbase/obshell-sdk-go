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

package model

import (
	"fmt"
	"time"
)

const (
	TASK_NAME_UPDATE_CONFIG    = "Update cluster config"
	TASK_NAME_UPDATE_OB_CONFIG = "Update observer config"

	DAG_INIT_CLUSTER   = "Initialize cluster"
	DAG_JOIN_TO_MASTER = "Join to master"
	DAG_JOIN_SELF      = "Join self"
)

type DagDetailDTO struct {
	*GenericDTO
	*DagDetail
}

type DagDetail struct {
	DagID    int64  `json:"dag_id" uri:"dag_id"`
	Name     string `json:"name"`
	Stage    int    `json:"stage"`
	MaxStage int    `json:"max_stage"`
	TaskStatusDTO
	AdditionalDataDTO
	Nodes []*NodeDetailDTO `json:"nodes"`
}

type TaskStatusDTO struct {
	State     string    `json:"state"`
	Operator  string    `json:"operator"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type AdditionalDataDTO struct {
	AdditionalData *map[string]any `json:"additional_data"`
}

type GenericDTO struct {
	GenericID string `json:"id" uri:"id" binding:"required"`
}

type NodeDetailDTO struct {
	*GenericDTO
	*NodeDetail
}

type NodeDetail struct {
	NodeID int64  `json:"node_id" uri:"node_id"`
	Name   string `json:"name"`
	TaskStatusDTO
	AdditionalDataDTO
	SubTasks []*TaskDetailDTO `json:"sub_tasks"`
}

type TaskDetailDTO struct {
	*GenericDTO
	*TaskDetail
}

type TaskDetail struct {
	TaskID int64  `json:"task_id" uri:"task_id"`
	Name   string `json:"name"`
	TaskStatusDTO
	AdditionalDataDTO
	ExecuteTimes int       `json:"execute_times"`
	ExecuteAgent AgentInfo `json:"execute_agent"`
	TaskLogs     []string  `json:"task_logs"`
}

func (t *TaskStatusDTO) IsRun() bool {
	return t.Operator == OPERATOR_MAP[RUN]
}

func (t *TaskStatusDTO) IsCancel() bool {
	return t.Operator == OPERATOR_MAP[CANCEL]
}

func (t *TaskStatusDTO) IsRetry() bool {
	return t.Operator == OPERATOR_MAP[RETRY]
}

func (t *TaskStatusDTO) IsRollback() bool {
	return t.Operator == OPERATOR_MAP[ROLLBACK]
}

func (t *TaskStatusDTO) IsPass() bool {
	return t.Operator == OPERATOR_MAP[PASS]
}

func (t *TaskStatusDTO) IsFailed() bool {
	return t.State == FAILED_STR
}

func (t *TaskStatusDTO) IsRunning() bool {
	return t.State == RUNNING_STR
}

func (t *TaskStatusDTO) IsSucceed() bool {
	return t.State == SUCCEED_STR
}

func (t *TaskStatusDTO) IsPending() bool {
	return t.State == PENDING_STR
}

func (t *TaskStatusDTO) IsReady() bool {
	return t.State == READY_STR
}

func (t *TaskStatusDTO) IsFinished() bool {
	return t.IsFailed() || t.IsSucceed()
}

func GetFailedDagLastLog(currentDag *DagDetailDTO) (res []string) {
	nodes := currentDag.Nodes

	var subTask *TaskDetailDTO
	var currentNode *NodeDetailDTO
	for i := 0; i < len(nodes); i++ {
		currentNode = nodes[i]
		if !currentNode.IsFailed() {
			continue
		}

		if currentNode.Operator == OPERATOR_MAP[CANCEL] {
			return append(res, fmt.Sprintf("Sorry, Task '%s' was cancelled", currentDag.Name))
		}

		for j := 0; j < len(currentNode.SubTasks); j++ {
			subTask = currentNode.SubTasks[j]
			if subTask.IsFailed() {
				lastLog := subTask.TaskLogs[len(subTask.TaskLogs)-1]
				res = append(res, fmt.Sprintf("%s %s", subTask.ExecuteAgent.String(), lastLog))
			}
		}
		return
	}
	return append(res, "No failed task log found, please check the task details")
}
