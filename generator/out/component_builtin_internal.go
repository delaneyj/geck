package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)


// Internal components are used by the ECS system and are not user-defined.
type Internal struct {
}

const (
    InternalID ID = 3
    InternalName = "Internal"
    InternalSizeBytes = 0
    InternalIsTag = true
    InternalIsBuiltin = true
)

var (
    InternalResetValue = Internal{
    }
    InternalByteOffsets = []int{
    }
    InternalIDSet = NewIDSet( InternalID )
)

func (c *Internal) Copy(other Internal) {
}

func (c *Internal) Reset() {
    c.Copy(InternalResetValue)
}

func (c *Internal) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, InternalSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Internal)(ptr)
    *tPtr = *c
    return buf
}

func (c *Internal) FromBytes(buf []byte) {
    if len(buf) < InternalSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Internal)(p)
	*c = *ptr
}

func (c *Internal) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Internal) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func InternalFromJSON(data []byte) (*Internal, error) {
    c := &Internal{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}



func (w *World) MarshalAllInternal() {

    log.Print("Marshaling Internal")
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
            if source == InternalID || target == InternalID {
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
        for i := range validCIDs {
            log.Printf("Marshaling Internal for %s", validNames[i])
            for _, e := range a.entities {
                log.Printf("%d", e)
                hasValues = true
            }
        }

    }

    if !hasValues {
        log.Print("No values for Internal")
    }
}

