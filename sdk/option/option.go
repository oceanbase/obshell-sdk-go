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

package option

type OptionType int

const (
	AUTH_OPT OptionType = iota + 1
)

type Optioner interface {
	Type() OptionType
	Name() string
	Value() interface{}
}

type BaseOption struct {
	name       string
	optionType OptionType
	value      interface{}
}

func NewBaseOption(name string, optionType OptionType, value interface{}) *BaseOption {
	return &BaseOption{
		name:       name,
		optionType: optionType,
		value:      value,
	}
}

func (o *BaseOption) Type() OptionType {
	return o.optionType
}

func (o *BaseOption) Name() string {
	return o.name
}

func (o *BaseOption) Value() interface{} {
	return o.value
}
