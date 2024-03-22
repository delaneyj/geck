package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)


// Allows instance-of relationships between entities
type InstanceOf struct {
    Value uint32 `json:"value,omitempty"`
}

const (
    InstanceOfID ID = 6
    InstanceOfName = "InstanceOf"
    InstanceOfSizeBytes = 4
    InstanceOfIsTag = false
    InstanceOfIsBuiltin = true
)

var (
    InstanceOfResetValue = InstanceOf{
        Value: 0,
    }
    InstanceOfByteOffsets = []int{
        0, // Value
    }
    InstanceOfIDSet = NewIDSet( InstanceOfID )
)

func (c *InstanceOf) Copy(other InstanceOf) {
    c.Value = other.Value
}

func (c *InstanceOf) Reset() {
    c.Copy(InstanceOfResetValue)
}

func (c *InstanceOf) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, InstanceOfSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*InstanceOf)(ptr)
    *tPtr = *c
    return buf
}

func (c *InstanceOf) FromBytes(buf []byte) {
    if len(buf) < InstanceOfSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*InstanceOf)(p)
	*c = *ptr
}

func (c *InstanceOf) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *InstanceOf) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func InstanceOfFromJSON(data []byte) (*InstanceOf, error) {
    c := &InstanceOf{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}



func (w *World) MarshalAllInstanceOf() {

    log.Print("Marshaling InstanceOf")
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
            if source == InstanceOfID || target == InstanceOfID {
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
            log.Printf("Marshaling InstanceOf for %s", validNames[i])

			colIdx, ok := colIndicies[cid]
			if !ok || colIdx < 0 {
				continue
			}

			col := a.dataColumns[colIdx]
            c := &InstanceOf{}

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
        log.Print("No values for InstanceOf")
    }
}

