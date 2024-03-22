package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)

type Velocity struct {
    X float32 `json:"x,omitempty"`
    Y float32 `json:"y,omitempty"`
    Z float32 `json:"z,omitempty"`
}

const (
    VelocityID ID = 1001
    VelocityName = "Velocity"
    VelocitySizeBytes = 12
    VelocityIsTag = false
    VelocityIsBuiltin = false
)

var (
    VelocityResetValue = Velocity{
        X: 0.000000,
        Y: 0.000000,
        Z: 0.000000,
    }
    VelocityByteOffsets = []int{
        0, // X
        4, // Y
        8, // Z
    }
    VelocityIDSet = NewIDSet( VelocityID )
)

func (c *Velocity) Copy(other Velocity) {
    c.X = other.X
    c.Y = other.Y
    c.Z = other.Z
}

func (c *Velocity) Reset() {
    c.Copy(VelocityResetValue)
}

func (c *Velocity) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, VelocitySizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Velocity)(ptr)
    *tPtr = *c
    return buf
}

func (c *Velocity) FromBytes(buf []byte) {
    if len(buf) < VelocitySizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Velocity)(p)
	*c = *ptr
}

func (c *Velocity) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Velocity) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func VelocityFromJSON(data []byte) (*Velocity, error) {
    c := &Velocity{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}

func (w *World) SetVelocitys(component Velocity, entities ...ID) {
    setComponentData(w, VelocityID, component, entities...)
}

func (w *World) Velocity(entity ID) *Velocity {
    data := &Velocity{}
    w.VelocityCopyTo(entity, data)
    return data
}

func (w *World) VelocityCopyTo(entity ID, copyTo *Velocity) {
    componentDataFromEntity(w, VelocityID, entity, copyTo)
}

func (w *World) RemoveVelocity(entities ...ID) {
    removeComponentFrom(w,  VelocityIDSet, NewIDSet(entities...))
}

func (w *World) MarshalAllVelocity() {

    log.Print("Marshaling Velocity")
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
            if source == VelocityID || target == VelocityID {
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
            log.Printf("Marshaling Velocity for %s", validNames[i])

			colIdx, ok := colIndicies[cid]
			if !ok || colIdx < 0 {
				continue
			}

			col := a.dataColumns[colIdx]
            c := &Velocity{}

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
        log.Print("No values for Velocity")
    }
}

