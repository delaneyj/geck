package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)


// A string identifier for an entity. This is used to identify an entity across the network and UI.
type Identifier struct {
    Value string `json:"value,omitempty"`
}

const (
    IdentifierID ID = 1
    IdentifierName = "Identifier"
    IdentifierSizeBytes = 16
    IdentifierIsTag = false
    IdentifierIsBuiltin = true
)

var (
    IdentifierResetValue = Identifier{
        Value: "!!!UNKNOWN IDENTIFIER!!!",
    }
    IdentifierByteOffsets = []int{
        0, // Value
    }
    IdentifierIDSet = NewIDSet( IdentifierID )
)

func (c *Identifier) Copy(other Identifier) {
    c.Value = other.Value
}

func (c *Identifier) Reset() {
    c.Copy(IdentifierResetValue)
}

func (c *Identifier) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, IdentifierSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Identifier)(ptr)
    *tPtr = *c
    return buf
}

func (c *Identifier) FromBytes(buf []byte) {
    if len(buf) < IdentifierSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Identifier)(p)
	*c = *ptr
}

func (c *Identifier) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Identifier) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func IdentifierFromJSON(data []byte) (*Identifier, error) {
    c := &Identifier{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}



func (w *World) MarshalAllIdentifier() {

    log.Print("Marshaling Identifier")
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
            if source == IdentifierID || target == IdentifierID {
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
            log.Printf("Marshaling Identifier for %s", validNames[i])

			colIdx, ok := colIndicies[cid]
			if !ok || colIdx < 0 {
				continue
			}

			col := a.dataColumns[colIdx]
            c := &Identifier{}

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
        log.Print("No values for Identifier")
    }
}

