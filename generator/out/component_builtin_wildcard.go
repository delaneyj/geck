package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)


// A wildcard tag. This is used to match any entity, used with pairs
type Wildcard struct {
}

const (
    WildcardID ID = 4
    WildcardName = "Wildcard"
    WildcardSizeBytes = 0
    WildcardIsTag = true
    WildcardIsBuiltin = true
)

var (
    WildcardResetValue = Wildcard{
    }
    WildcardByteOffsets = []int{
    }
    WildcardIDSet = NewIDSet( WildcardID )
)

func (c *Wildcard) Copy(other Wildcard) {
}

func (c *Wildcard) Reset() {
    c.Copy(WildcardResetValue)
}

func (c *Wildcard) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, WildcardSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Wildcard)(ptr)
    *tPtr = *c
    return buf
}

func (c *Wildcard) FromBytes(buf []byte) {
    if len(buf) < WildcardSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Wildcard)(p)
	*c = *ptr
}

func (c *Wildcard) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Wildcard) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func WildcardFromJSON(data []byte) (*Wildcard, error) {
    c := &Wildcard{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}



func (w *World) MarshalAllWildcard() {

    log.Print("Marshaling Wildcard")
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
            if source == WildcardID || target == WildcardID {
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
            log.Printf("Marshaling Wildcard for %s", validNames[i])
            for _, e := range a.entities {
                log.Printf("%d", e)
                hasValues = true
            }
        }

    }

    if !hasValues {
        log.Print("No values for Wildcard")
    }
}

