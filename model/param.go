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

const (
	RUN = iota + 1
	RETRY
	ROLLBACK
	CANCEL
	PASS

	RUN_STR      = "RUN"
	RETRY_STR    = "RETRY"
	ROLLBACK_STR = "ROLLBACK"
	CANCEL_STR   = "CANCEL"
	PASS_STR     = "PASS"
)

// State
const (
	PENDING = iota + 1
	READY
	RUNNING
	FAILED
	SUCCEED

	PENDING_STR = "PENDING"
	READY_STR   = "READY"
	RUNNING_STR = "RUNNING"
	FAILED_STR  = "FAILED"
	SUCCEED_STR = "SUCCEED"
)

var STATE_MAP = map[int]string{
	PENDING: PENDING_STR,
	READY:   READY_STR,
	RUNNING: RUNNING_STR,
	FAILED:  FAILED_STR,
	SUCCEED: SUCCEED_STR,
}

var OPERATOR_MAP = map[int]string{
	RUN:      RUN_STR,
	RETRY:    RETRY_STR,
	ROLLBACK: ROLLBACK_STR,
	CANCEL:   CANCEL_STR,
	PASS:     PASS_STR,
}
