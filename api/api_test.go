package api

import (
	"github.com/B1NARY-GR0UP/openalysis/config"
	"testing"
)

func TestAPI(t *testing.T) {
	AddGroups(config.Group{Name: "hello"}, config.Group{Name: "world"})
}
