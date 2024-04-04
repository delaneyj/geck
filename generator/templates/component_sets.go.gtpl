package {{.PackageName}}

type {{.Name.Singular.Pascal}} struct {
	lastIdx int

	// owned components
    {{range .OwnedComponents -}}
	owned{{.Name.Plural.Pascal}}Store    *SparseSet[{{.Name.Singular.Pascal}}]
    {{end -}}
    {{ if .BorrowedComponents -}}
    borrowedSparseSet *SparseSet[Entity]
    {{range .BorrowedComponents -}}
    borrowed{{.Name.Plural.Pascal}}Store *SparseSet[{{.Name.Singular.Pascal}}]
    {{end -}}
    {{end -}}
}

func New{{.Name.Singular.Pascal}}(w *World) *{{.Name.Singular.Pascal}} {
    set := &{{.Name.Singular.Pascal}}{
        lastIdx: -1,

        {{range .OwnedComponents -}}
        owned{{.Name.Plural.Pascal}}Store : w.{{.Name.Plural.Camel}}Store,
        {{end -}}
        {{ if .BorrowedComponents -}}
        borrowedSparseSet : NewSparseSet[Entity](),
        {{range .BorrowedComponents -}}
        borrowed{{.Name.Plural.Pascal}}Store : NewSparseSet[{{.Name.Singular.Pascal}}](),
        {{end -}}
        {{end -}}
    }
    return set
}

func (set *{{.Name.Singular.Pascal}}) PossibleUpdate(entities ...Entity) {
    for _, e := range entities {
        hasAllOwned := true
        {{range .OwnedComponents}}
        if !set.owned{{.Name.Plural.Pascal}}Store.Has(e) {
            hasAllOwned = false
            break
        }
        {{end}}

        sparseIdx := e.Index()

        if hasAllOwned {
            // swap with next after last
            set.lastIdx++

            wasSwapped := false
            {{range .OwnedComponents }}
            if set.owned{{.Name.Plural.Pascal}}Store.ownedSetSwap(set.lastIdx, sparseIdx, false) {
                wasSwapped = true
            }
            {{end }}

            if !wasSwapped {
                set.lastIdx--
            }
        } else {
            // swap with last
            wasSwapped := false
            {{range .OwnedComponents}}
            if set.owned{{.Name.Plural.Pascal}}Store.ownedSetSwap(set.lastIdx, sparseIdx, true) {
                wasSwapped = true
            }
            {{end}}

            if wasSwapped {
                set.lastIdx--
            }
        }

        // do something with
        // hasAllBorrowed := true


    }
}

func (set *{{.Name.Singular.Pascal}}) Len() int {
    return set.lastIdx + 1
}

func (set *{{.Name.Singular.Pascal}}) NewIterator() *{{.Name.Singular.Pascal}}Iter {
    iter := &{{.Name.Singular.Pascal}}Iter{ set: set }
    iter.Reset()
    return iter
}

type {{.Name.Singular.Pascal}}Iter struct {
    set *{{.Name.Singular.Pascal}}
    current int
}

func (iter *{{.Name.Singular.Pascal}}Iter) Reset() {
    {{ if .HasWriteableComponents -}}
    iter.current = iter.set.lastIdx
    {{else -}}
	iter.current = 0
    {{end -}}
}

func (iter *{{.Name.Singular.Pascal}}Iter) HasNext() bool {
    {{if .HasWriteableComponents -}}
    return iter.current >= 0
    {{else -}}
	return iter.current <= iter.set.lastIdx
    {{end -}}
}

func (iter *{{.Name.Singular.Pascal}}Iter) Next() (
    Entity,
    {{range .OwnedComponents -}}
    {{if .IsWritable}}*{{end}}{{.Name.Singular.Pascal}},
    {{end -}}
    {{range .BorrowedComponents -}}
    {{if .IsWritable}}*{{end}}{{.Name.Singular.Pascal}},
    {{end -}}
) {
    {{range $i, $c := .OwnedComponents -}}
    {{ $v := printf "comp%x" $i -}}
    {{if eq $i 0 -}}
    e := iter.set.owned{{$c.Name.Plural.Pascal}}Store.dense[iter.current]
    {{end -}}
    {{$v}} := {{if .IsWritable}}&{{end}}iter.set.owned{{$c.Name.Plural.Pascal}}Store.components[iter.current]
    {{end -}}

    {{if .HasWriteableComponents -}}
    iter.current--
    {{else -}}
    iter.current++
    {{end -}}

    return e, {{range $i, $c := .OwnedComponents}}{{if $i}},{{end}}{{$v := printf "comp%x" $i}}{{$v}}{{end}}
}


