package tests

import (
	"testing"

	"github.com/delaneyj/geck"
)

func TestEntity(t *testing.T) {
	var table = []struct {
		index      uint32
		generation uint16
		value      geck.ID
	}{
		{0, 0, 0x0000000000000000},
		{0, 1, 0x0000000100000000},
		{1, 0, 0x0000000000000001},
		{1, 1, 0x0000000100000001},
		{2, 0, 0x0000000000000002},
		{2, 1, 0x0000000100000002},
		{3, 0, 0x0000000000000003},
		{3, 1, 0x0000000100000003},
		{123456789, 0, 0x00000000075BCD15},
		{123456789, 1, 0x00000001075BCD15},
	}

	for _, row := range table {
		e := row.value
		if e.Index() != row.index {
			t.Errorf("Expected index to be 2, got %d", e.Index())
		}
		if e.Generation() != row.generation {
			t.Errorf("Expected generation to be 1, got %d", e.Generation())
		}
	}
}
