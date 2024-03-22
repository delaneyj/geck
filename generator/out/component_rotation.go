package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)

type Rotation struct {
    X float32 `json:"x,omitempty"`
    Y float32 `json:"y,omitempty"`
    Z float32 `json:"z,omitempty"`
    W float32 `json:"w,omitempty"`
}

const (
    RotationID ID = 1002
    RotationName = "Rotation"
    RotationSizeBytes = 16
    RotationIsTag = false
    RotationIsBuiltin = false
)

var (
    RotationResetValue = Rotation{
        X: 0.000000,
        Y: 0.000000,
        Z: 0.000000,
        W: 1.000000,
    }
    RotationByteOffsets = []int{
        0, // X
        4, // Y
        8, // Z
        12, // W
    }
    RotationIDSet = NewIDSet( RotationID )
)

func (c *Rotation) Copy(other Rotation) {
    c.X = other.X
    c.Y = other.Y
    c.Z = other.Z
    c.W = other.W
}

func (c *Rotation) Reset() {
    c.Copy(RotationResetValue)
}

func (c *Rotation) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, RotationSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Rotation)(ptr)
    *tPtr = *c
    return buf
}

func (c *Rotation) FromBytes(buf []byte) {
    if len(buf) < RotationSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Rotation)(p)
	*c = *ptr
}

func (c *Rotation) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Rotation) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func RotationFromJSON(data []byte) (*Rotation, error) {
    c := &Rotation{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}

func (w *World) SetRotations(component Rotation, entities ...ID) {
    setComponentData(w, RotationID, component, entities...)
}

func (w *World) Rotation(entity ID) *Rotation {
    data := &Rotation{}
    w.RotationCopyTo(entity, data)
    return data
}

func (w *World) RotationCopyTo(entity ID, copyTo *Rotation) {
    componentDataFromEntity(w, RotationID, entity, copyTo)
}

func (w *World) RemoveRotation(entities ...ID) {
    removeComponentFrom(w,  RotationIDSet, NewIDSet(entities...))
}

func (w *World) MarshalAllRotation() {

    log.Print("Marshaling Rotation")
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
            if source == RotationID || target == RotationID {
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
            log.Printf("Marshaling Rotation for %s", validNames[i])

			colIdx, ok := colIndicies[cid]
			if !ok || colIdx < 0 {
				continue
			}

			col := a.dataColumns[colIdx]
            c := &Rotation{}

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
        log.Print("No values for Rotation")
    }
}

