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
	"errors"

	"github.com/oceanbase/obshell-sdk-go/model"
)

// Aggregation Functions
// Clear clears the agent status to SINGLE before the the cluster init successfully.
func (c *Client) Clear() (err error) {
	agentStatus, err := c.GetStatus()
	if err != nil {
		return err
	}

	lastDag, err := c.GetAgentLastMaintenanceDag()
	if err != nil {
		return err
	}

	need_remove := false
	if agentStatus.Agent.Identity == model.FOLLOWER || agentStatus.Agent.Identity == model.MASTER {
		need_remove = true
	}

	// init rollback
	if lastDag.Name == model.DAG_INIT_CLUSTER {
		if !lastDag.IsFinished() {
			return errors.New("the 'Initialize Cluster' task is not finished yet")
		}

		if lastDag.IsSucceed() && lastDag.IsRun() {
			return errors.New("the 'Initialize Cluster' task is already succeeded")
		}

		if lastDag.IsFailed() {
			rollbackRequest := c.NewOperateDagRequest(lastDag.GenericID, model.ROLLBACK_STR)
			if err := c.OperateDagSyncWithRequest(rollbackRequest); err != nil {
				return err
			}
		}

		need_remove = true
	}

	if need_remove {
		removeRequest := c.NewRemoveRequest(c.GetHost(), c.GetPort())
		if _, err := c.RemoveSyncWithRequest(removeRequest); err != nil {
			return err
		}
	}

	return nil
}
