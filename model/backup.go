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

import "time"

type CdbObBackupTask struct {
	TenantID              int64     `json:"tenant_id"`
	TaskID                int64     `json:"task_id"`
	JobID                 int64     `json:"job_id"`
	Incarnation           int64     `json:"incarnation"`
	BackupSetID           int64     `json:"backup_set_id"`
	StartTimestamp        time.Time `json:"start_timestamp"`
	EndTimestamp          time.Time `json:"end_timestamp"`
	Status                string    `json:"status"`
	StartScn              int64     `json:"start_scn"`
	EndScn                int64     `json:"end_scn"`
	UserLsStartScn        int64     `json:"user_ls_start_scn"`
	EncryptionMode        string    `json:"encryption_mode"`
	InputBytes            int64     `json:"input_bytes"`
	OutputBytes           int64     `json:"output_bytes"`
	OutputRateBytes       float64   `json:"output_rate_bytes"`
	ExtraMetaBytes        int64     `json:"extra_meta_bytes"`
	TabletCount           int64     `json:"tablet_count"`
	FinishTabletCount     int64     `json:"finish_tablet_count"`
	MacroBlockCount       int64     `json:"macro_block_count"`
	FinishMacroBlockCount int64     `json:"finish_macro_block_count"`
	FileCount             int64     `json:"file_count"`
	MetaTurnID            int64     `json:"meta_turn_id"`
	DataTurnID            int64     `json:"data_turn_id"`
	Result                int64     `json:"result"`
	Comment               string    `json:"comment"`
	Path                  string    `json:"path"`
}
