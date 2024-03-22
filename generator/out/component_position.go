package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)

type Position struct {
    X float32 `json:"x,omitempty"`
    Y float32 `json:"y,omitempty"`
    Z float32 `json:"z,omitempty"`
}

const (
    PositionID ID = 1000
    PositionName = "Position"
    PositionSizeBytes = 12
    PositionIsTag = false
    PositionIsBuiltin = false
)

var (
    PositionResetValue = Position{
        X: 0.000000,
        Y: 0.000000,
        Z: 0.000000,
    }
    PositionByteOffsets = []int{
        0, // X
        4, // Y
        8, // Z
    }
    PositionIDSet = NewIDSet( PositionID )
)

func (c *Position) Copy(other Position) {
    c.X = other.X
    c.Y = other.Y
    c.Z = other.Z
}

func (c *Position) Reset() {
    c.Copy(PositionResetValue)
}

func (c *Position) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, PositionSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Position)(ptr)
    *tPtr = *c
    return buf
}

func (c *Position) FromBytes(buf []byte) {
    if len(buf) < PositionSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Position)(p)
	*c = *ptr
}

func (c *Position) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Position) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func PositionFromJSON(data []byte) (*Position, error) {
    c := &Position{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}

func (w *World) SetPositions(component Position, entities ...ID) {
    setComponentData(w, PositionID, component, entities...)
}

func (w *World) Position(entity ID) *Position {
    data := &Position{}
    w.PositionCopyTo(entity, data)
    return data
}

func (w *World) PositionCopyTo(entity ID, copyTo *Position) {
    componentDataFromEntity(w, PositionID, entity, copyTo)
}

func (w *World) RemovePosition(entities ...ID) {
    removeComponentFrom(w,  PositionIDSet, NewIDSet(entities...))
}

func (w *World) MarshalAllPosition() {

    log.Print("Marshaling Position")
    hasValues := false
    for _, a := range w.archetypes {
        if len(a.entities) == 0 {
			continue
		}

		count := a.componentIDs.Cardinality()
		if count == 0 {
			continue
		}

		validCIDs := []ID{}
        validNames := []string{}
		a.componentIDs.Range(func(cid ID) {
			source, target, _ := cid.SplitPair()
            if source == PositionID || target == PositionID {
				validCIDs = append(validCIDs, cid)

                sn := w.EntityName(source)
                if sn == "" {
                    sn = "_"
                }
                tn := w.EntityName(target)
                if tn == "" {
                    tn = "_"
                }
                n := fmt.Sprintf("%s,%s", sn, tn)
                validNames = append(validNames, n)
			}

		})
		if len(validCIDs) == 0 {
			continue
		}
		colIndicies := w.archetypeComponentColumnIndicies[a.hash]
		for i, cid := range validCIDs {
            log.Printf("Marshaling Position for %s", validNames[i])

			colIdx, ok := colIndicies[cid]
			if !ok || colIdx < 0 {
				continue
			}

			col := a.dataColumns[colIdx]
            c := &Position{}

			for i, e := range a.entities {
				start := uintptr(i) * col.metadata.elementSize
				end := start + col.metadata.elementSize
				buf := col.data[start:end]
				c.FromBytes(buf)

                log.Printf("%d : %+v", e,c)
                hasValues = true
			}
		}

    }

    if !hasValues {
        log.Print("No values for Position")
    }
}

