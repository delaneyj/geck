package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)


// Allows parent-child relationships between entities
type ChildOf struct {
    Parent uint32 `json:"parent,omitempty"`
}

const (
    ChildOfID ID = 5
    ChildOfName = "ChildOf"
    ChildOfSizeBytes = 4
    ChildOfIsTag = false
    ChildOfIsBuiltin = true
)

var (
    ChildOfResetValue = ChildOf{
        Parent: 0,
    }
    ChildOfByteOffsets = []int{
        0, // Parent
    }
    ChildOfIDSet = NewIDSet( ChildOfID )
)

func (c *ChildOf) Copy(other ChildOf) {
    c.Parent = other.Parent
}

func (c *ChildOf) Reset() {
    c.Copy(ChildOfResetValue)
}

func (c *ChildOf) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, ChildOfSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*ChildOf)(ptr)
    *tPtr = *c
    return buf
}

func (c *ChildOf) FromBytes(buf []byte) {
    if len(buf) < ChildOfSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*ChildOf)(p)
	*c = *ptr
}

func (c *ChildOf) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *ChildOf) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func ChildOfFromJSON(data []byte) (*ChildOf, error) {
    c := &ChildOf{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}



func (w *World) MarshalAllChildOf() {

    log.Print("Marshaling ChildOf")
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
            if source == ChildOfID || target == ChildOfID {
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
            log.Printf("Marshaling ChildOf for %s", validNames[i])

			colIdx, ok := colIndicies[cid]
			if !ok || colIdx < 0 {
				continue
			}

			col := a.dataColumns[colIdx]
            c := &ChildOf{}

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
        log.Print("No values for ChildOf")
    }
}

