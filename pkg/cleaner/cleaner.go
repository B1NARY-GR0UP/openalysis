// Copyright 2024 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cleaner

import "strings"

type Cleaner struct {
	Strategies map[string]string
}

// New a Cleaner which accept strategies in "`Before` => `After`" format
// e.g. "`JustLorain` => `justlorain`"
func New(strategies ...string) *Cleaner {
	cleaner := &Cleaner{
		Strategies: make(map[string]string, len(strategies)),
	}
	for _, strategy := range strategies {
		parts := strings.Split(strategy, "=>")
		k := parts[0][strings.Index(parts[0], "`")+1 : strings.LastIndex(parts[0], "`")]
		v := parts[1][strings.Index(parts[1], "`")+1 : strings.LastIndex(parts[1], "`")]
		cleaner.Strategies[k] = v
	}
	return cleaner
}

func (c *Cleaner) Clean(input string) string {
	if v, ok := c.Strategies[input]; ok {
		return v
	}
	return input
}
