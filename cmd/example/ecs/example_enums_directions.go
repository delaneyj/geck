package ecs

import "fmt"

type EnumDirection uint32

const (
	EnumDirectionUnknown EnumDirection = 0
	EnumDirectionNorth   EnumDirection = 1
	EnumDirectionSouth   EnumDirection = 2
	EnumDirectionEast    EnumDirection = 4
	EnumDirectionWest    EnumDirection = 8
)

func EnumDirectionFromString(value string) EnumDirection {
	switch value {
	case "unknown":
		return EnumDirectionUnknown
	case "north":
		return EnumDirectionNorth
	case "south":
		return EnumDirectionSouth
	case "east":
		return EnumDirectionEast
	case "west":
		return EnumDirectionWest
	default:
		panic(fmt.Sprintf("Unknown value for EnumDirection: %s", value))
	}
}

func (e EnumDirection) String() string {
	switch e {
	case EnumDirectionUnknown:
		return "unknown"
	case EnumDirectionNorth:
		return "north"
	case EnumDirectionSouth:
		return "south"
	case EnumDirectionEast:
		return "east"
	case EnumDirectionWest:
		return "west"
	default:
		panic(fmt.Sprintf("Unknown value for EnumDirection: %d", e))
	}
}

func (e EnumDirection) U32() uint32 {
	return uint32(e)
}
