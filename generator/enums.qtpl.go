// Code generated by qtc from "enums.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// package generator
//

//line generator/enums.qtpl:3
package generator

//line generator/enums.qtpl:3
import "fmt"

//line generator/enums.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line generator/enums.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line generator/enums.qtpl:5
func streamenumTemplate(qw422016 *qt422016.Writer, data *enumTmplData) {
//line generator/enums.qtpl:5
	qw422016.N().S(`
package `)
//line generator/enums.qtpl:6
	qw422016.E().S(data.PackageName)
//line generator/enums.qtpl:6
	qw422016.N().S(`
    
`)
//line generator/enums.qtpl:9
	enumName := fmt.Sprintf("Enum%s", data.Name.Singular.Pascal)

//line generator/enums.qtpl:10
	qw422016.N().S(`
type `)
//line generator/enums.qtpl:11
	qw422016.E().S(enumName)
//line generator/enums.qtpl:11
	qw422016.N().S(` uint32

const (
`)
//line generator/enums.qtpl:14
	for _, value := range data.Values {
//line generator/enums.qtpl:14
		qw422016.N().S(`    `)
//line generator/enums.qtpl:15
		qw422016.E().S(enumName)
//line generator/enums.qtpl:15
		qw422016.E().S(value.Name.Singular.Pascal)
//line generator/enums.qtpl:15
		qw422016.N().S(` `)
//line generator/enums.qtpl:15
		qw422016.E().S(enumName)
//line generator/enums.qtpl:15
		qw422016.N().S(` = `)
//line generator/enums.qtpl:15
		qw422016.N().D(value.Value)
//line generator/enums.qtpl:15
		qw422016.N().S(`
`)
//line generator/enums.qtpl:16
	}
//line generator/enums.qtpl:16
	qw422016.N().S(`)

func `)
//line generator/enums.qtpl:19
	qw422016.E().S(enumName)
//line generator/enums.qtpl:19
	qw422016.N().S(`FromString(value string) `)
//line generator/enums.qtpl:19
	qw422016.E().S(enumName)
//line generator/enums.qtpl:19
	qw422016.N().S(` {
    switch value {
`)
//line generator/enums.qtpl:21
	for _, value := range data.Values {
//line generator/enums.qtpl:21
		qw422016.N().S(`    case "`)
//line generator/enums.qtpl:22
		qw422016.E().S(value.Name.Singular.Snake)
//line generator/enums.qtpl:22
		qw422016.N().S(`":
        return `)
//line generator/enums.qtpl:23
		qw422016.E().S(enumName)
//line generator/enums.qtpl:23
		qw422016.E().S(value.Name.Singular.Pascal)
//line generator/enums.qtpl:23
		qw422016.N().S(`
`)
//line generator/enums.qtpl:24
	}
//line generator/enums.qtpl:24
	qw422016.N().S(`    default:
        panic(fmt.Sprintf("Unknown value for `)
//line generator/enums.qtpl:26
	qw422016.E().S(enumName)
//line generator/enums.qtpl:26
	qw422016.N().S(`: %s", value))
    }
}

func (e `)
//line generator/enums.qtpl:30
	qw422016.E().S(enumName)
//line generator/enums.qtpl:30
	qw422016.N().S(`) String() string {
    switch e {
`)
//line generator/enums.qtpl:32
	for _, value := range data.Values {
//line generator/enums.qtpl:32
		qw422016.N().S(`    case `)
//line generator/enums.qtpl:33
		qw422016.E().S(enumName)
//line generator/enums.qtpl:33
		qw422016.E().S(value.Name.Singular.Pascal)
//line generator/enums.qtpl:33
		qw422016.N().S(`:
        return "`)
//line generator/enums.qtpl:34
		qw422016.E().S(value.Name.Singular.Snake)
//line generator/enums.qtpl:34
		qw422016.N().S(`"
`)
//line generator/enums.qtpl:35
	}
//line generator/enums.qtpl:35
	qw422016.N().S(`    default:
        panic(fmt.Sprintf("Unknown value for `)
//line generator/enums.qtpl:37
	qw422016.E().S(enumName)
//line generator/enums.qtpl:37
	qw422016.N().S(`: %d", e))
    }
}

func (e `)
//line generator/enums.qtpl:41
	qw422016.E().S(enumName)
//line generator/enums.qtpl:41
	qw422016.N().S(`) U32() uint32 {
    return uint32(e)
}

`)
//line generator/enums.qtpl:45
}

//line generator/enums.qtpl:45
func writeenumTemplate(qq422016 qtio422016.Writer, data *enumTmplData) {
//line generator/enums.qtpl:45
	qw422016 := qt422016.AcquireWriter(qq422016)
//line generator/enums.qtpl:45
	streamenumTemplate(qw422016, data)
//line generator/enums.qtpl:45
	qt422016.ReleaseWriter(qw422016)
//line generator/enums.qtpl:45
}

//line generator/enums.qtpl:45
func enumTemplate(data *enumTmplData) string {
//line generator/enums.qtpl:45
	qb422016 := qt422016.AcquireByteBuffer()
//line generator/enums.qtpl:45
	writeenumTemplate(qb422016, data)
//line generator/enums.qtpl:45
	qs422016 := string(qb422016.B)
//line generator/enums.qtpl:45
	qt422016.ReleaseByteBuffer(qb422016)
//line generator/enums.qtpl:45
	return qs422016
//line generator/enums.qtpl:45
}