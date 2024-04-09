package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type EnumDirection int

const (
	EnumDirectionDirectionUnknown = 0
	EnumDirectionNorth            = 1
	EnumDirectionSouth            = 2
	EnumDirectionEast             = 4
	EnumDirectionWest             = 8
)

func (e EnumDirection) String() (string, bool) {
	switch e {
	case EnumDirectionDirectionUnknown:
		return "DirectionUnknown", true
	case EnumDirectionNorth:
		return "North", true
	case EnumDirectionSouth:
		return "South", true
	case EnumDirectionEast:
		return "East", true
	case EnumDirectionWest:
		return "West", true

	default:
		return "", false
	}
}

func (e EnumDirection) ToInt() int {
	return int(e)
}

func EnumDirectionFromInt(i int) EnumDirection {
	return EnumDirection(i)
}

func (e EnumDirection) ToPB() ecspb.DirectionEnum {
	return ecspb.DirectionEnum(e.ToInt())
}

func EnumDirectionSet(flags ...EnumDirection) EnumDirection {
	var e EnumDirection
	for _, flag := range flags {
		e |= flag
	}
	return e
}

func (e EnumDirection) Has(flags ...EnumDirection) bool {
	for _, flag := range flags {
		if e&flag == 0 {
			return false
		}
	}
	return true
}

func (e EnumDirection) Set(flags ...EnumDirection) EnumDirection {
	for _, flag := range flags {
		e |= flag
	}
	return e
}

func (e EnumDirection) Clear(flags ...EnumDirection) EnumDirection {
	for _, flag := range flags {
		e &= ^flag
	}
	return e
}

func (e EnumDirection) Toggle(flags ...EnumDirection) EnumDirection {
	for _, flag := range flags {
		e ^= flag
	}
	return e
}

func (e EnumDirection) ToggleAll() EnumDirection {
	return e ^ EnumDirectionSet(EnumDirectionDirectionUnknown,
		EnumDirectionNorth,
		EnumDirectionSouth,
		EnumDirectionEast,
		EnumDirectionWest,
	)
}

func (e EnumDirection) AllSet() (flags []EnumDirection) {

	if e&EnumDirectionDirectionUnknown != 0 {
		flags = append(flags, EnumDirectionDirectionUnknown)
	}
	if e&EnumDirectionNorth != 0 {
		flags = append(flags, EnumDirectionNorth)
	}
	if e&EnumDirectionSouth != 0 {
		flags = append(flags, EnumDirectionSouth)
	}
	if e&EnumDirectionEast != 0 {
		flags = append(flags, EnumDirectionEast)
	}
	if e&EnumDirectionWest != 0 {
		flags = append(flags, EnumDirectionWest)
	}

	return flags
}
