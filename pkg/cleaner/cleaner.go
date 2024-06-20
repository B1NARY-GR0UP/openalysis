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

import (
	"fmt"
	"strings"
	"sync"
)

type Cleaner struct {
	mu         sync.Mutex
	Strategies map[string]string
}

// New a Cleaner
func New() *Cleaner {
	return &Cleaner{
		Strategies: make(map[string]string),
	}
}

// AddStrategies accepts strategies in "`Before` => `After`" format
// e.g. "`JustLorain` => `justlorain`"
func (c *Cleaner) AddStrategies(strategies ...string) error {
	for _, strategy := range strategies {
		parts := strings.Split(strategy, "=>")
		if len(parts) != 2 {
			return fmt.Errorf("invalid cleaner strategy format: %s", strategy)
		}

		k := strings.Trim(parts[0], "` ")
		v := strings.Trim(parts[1], "` ")

		c.mu.Lock()
		c.Strategies[k] = v
		c.mu.Unlock()
	}
	return nil
}

func (c *Cleaner) DeleteStrategies(before string) {
	c.mu.Lock()
	delete(c.Strategies, before)
	c.mu.Unlock()
}

func (c *Cleaner) Clean(input string) string {
	if v, ok := c.Strategies[input]; ok {
		return v
	}
	return input
}
