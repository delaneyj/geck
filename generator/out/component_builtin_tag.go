package out

import (
    "github.com/goccy/go-json"
    "unsafe"
    "log"
    "fmt"
)


// A tag
type Tag struct {
}

const (
    TagID ID = 8
    TagName = "Tag"
    TagSizeBytes = 0
    TagIsTag = true
    TagIsBuiltin = true
)

var (
    TagResetValue = Tag{
    }
    TagByteOffsets = []int{
    }
    TagIDSet = NewIDSet( TagID )
)

func (c *Tag) Copy(other Tag) {
}

func (c *Tag) Reset() {
    c.Copy(TagResetValue)
}

func (c *Tag) Bytes() []byte {
    // copy the data into the pointer using unsafe
    buf := make([]byte, TagSizeBytes)
	ptr := unsafe.Pointer(&buf[0])
    tPtr := (*Tag)(ptr)
    *tPtr = *c
    return buf
}

func (c *Tag) FromBytes(buf []byte) {
    if len(buf) < TagSizeBytes {
        panic("Invalid buffer size")
    }
    p := unsafe.Pointer(unsafe.SliceData(buf))
    ptr := (*Tag)(p)
	*c = *ptr
}

func (c *Tag) ToJSON() ([]byte, error) {
    return json.Marshal(c)
}

func (c *Tag) FromJSON(data []byte) error {
    return json.Unmarshal(data, c)
}

func TagFromJSON(data []byte) (*Tag, error) {
    c := &Tag{}
    if err := c.FromJSON(data); err != nil {
        return nil, err
    }
    return c, nil
}



func (w *World) MarshalAllTag() {

    log.Print("Marshaling Tag")
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
            if source == TagID || target == TagID {
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
            log.Printf("Marshaling Tag for %s", validNames[i])
            for _, e := range a.entities {
                log.Printf("%d", e)
                hasValues = true
            }
        }

    }

    if !hasValues {
        log.Print("No values for Tag")
    }
}

