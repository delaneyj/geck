package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)


// A component tag
type Component struct {
}

const (
    ComponentID ID = 7
    ComponentName = "Component"
    ComponentSizeBytes = 0
    ComponentIsTag = true
    ComponentIsBuiltin = true
)

var (
    ComponentResetValue = Component{
    }
    ComponentByteOffsets = []int{
    }
    ComponentIDSet = NewIDSet( ComponentID )
)

func (c *Component) Copy(other Component) {
}

func (c *Component) Reset() {
    c.Copy(ComponentResetValue)
}

func (c *Component) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, ComponentSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Component)(ptr)
    *tPtr = *c
    return buf
}

func (c *Component) FromBytes(buf []byte) {
    if len(buf) < ComponentSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Component)(p)
	*c = *ptr
}

func (c *Component) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Component) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func ComponentFromJSON(data []byte) (*Component, error) {
    c := &Component{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}



func (w *World) MarshalAllComponent() {

    log.Print("Marshaling Component")
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
            if source == ComponentID || target == ComponentID {
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
            log.Printf("Marshaling Component for %s", validNames[i])
            for _, e := range a.entities {
                log.Printf("%d", e)
                hasValues = true
            }
        }

    }

    if !hasValues {
        log.Print("No values for Component")
    }
}

