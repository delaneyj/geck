package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)


// A name tag
type Name struct {
}

const (
    NameID ID = 2
    NameName = "Name"
    NameSizeBytes = 0
    NameIsTag = true
    NameIsBuiltin = true
)

var (
    NameResetValue = Name{
    }
    NameByteOffsets = []int{
    }
    NameIDSet = NewIDSet( NameID )
)

func (c *Name) Copy(other Name) {
}

func (c *Name) Reset() {
    c.Copy(NameResetValue)
}

func (c *Name) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, NameSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Name)(ptr)
    *tPtr = *c
    return buf
}

func (c *Name) FromBytes(buf []byte) {
    if len(buf) < NameSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Name)(p)
	*c = *ptr
}

func (c *Name) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Name) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func NameFromJSON(data []byte) (*Name, error) {
    c := &Name{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}



func (w *World) MarshalAllName() {

    log.Print("Marshaling Name")
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
            if source == NameID || target == NameID {
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
            log.Printf("Marshaling Name for %s", validNames[i])
            for _, e := range a.entities {
                log.Printf("%d", e)
                hasValues = true
            }
        }

    }

    if !hasValues {
        log.Print("No values for Name")
    }
}

