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

type ResourcePoolInfo struct {
	Name         string `json:"name"`
	Id           int    `json:"id"`
	ZoneList     string `json:"zone_list"`
	UnitNum      int    `json:"unit_num"`
	UnitConfigId int    `json:"unit_config_id"`
	TenantId     int    `json:"tenant_id"`
}

type TenantOverview struct {
	Name         string    `json:"tenant_name"`
	Id           int       `json:"tenant_id"`
	CreatedTime  time.Time `json:"created_time"`
	Mode         string    `json:"mode"`
	Status       string    `json:"status"`
	Locked       string    `json:"locked"`
	PrimaryZone  string    `json:"primary_zone"`
	Locality     string    `json:"locality"`
	InRecyclebin string    `json:"in_recyclebin"`
}

type TenantInfo struct {
	TenantOverview
	Charset   string                  `json:"charset"`   // Only for ORACLE tenant
	Collation string                  `json:"collation"` // Only for ORACLE tenant
	WhiteList string                  `json:"whitelist"`
	Pools     []*ResourcePoolWithUnit `json:"pools"`
}

type ResourcePoolWithUnit struct {
	Name     string              `json:"pool_name"`
	Id       int                 `json:"pool_id"`
	ZoneList string              `json:"zone_list"`
	UnitNum  int                 `json:"unit_num"`
	Unit     *ResourceUnitConfig `json:"unit_config"`
}

type ResourceUnitConfig struct {
	GmtCreate    time.Time `json:"create_time"`
	GmtModified  time.Time `json:"modify_time"`
	UnitConfigId int       `json:"unit_config_id"`
	Name         string    `json:"name"`
	MaxCpu       float64   `json:"max_cpu"`
	MinCpu       float64   `json:"min_cpu"`
	MemorySize   int       `json:"memory_size"`
	LogDiskSize  int       `json:"log_disk_size"`
	MaxIops      int       `json:"max_iops"`
	MinIops      int       `json:"min_iops"`
}

type VariableInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Info  string `json:"info"`
}

type ParameterInfo struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	DataType  string `json:"data_type"`
	Info      string `json:"info"`
	EditLevel string `json:"edit_level"`
}

type RecycledTenantOverView struct {
	Name         string `json:"object_name"`
	OriginalName string `json:"original_tenant_name"`
	CanUndrop    string `json:"can_undrop"`
	CanPurge     string `json:"can_purge"`
}
