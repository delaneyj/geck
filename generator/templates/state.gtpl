package {{.PackageName}}

import (
    "fmt"
    "sync"
    "github.com/google/uuid"
)

type {{.Name.Pascal}}State int

const (
    {{range $i, $s := .States -}}
    {{$.Name.Pascal}}{{$s.Pascal}} {{$.Name.Pascal}}State = {{$i}}
    {{end -}}
)

type {{.Name.Pascal}}UpdateFunc func(s {{.Name.Pascal}}State, isPaused bool) error

type {{.Name.Pascal}}StateInstance struct {
    current {{.Name.Pascal}}State
    isPaused bool
    mu *sync.RWMutex
    subs map[uuid.UUID]{{.Name.Pascal}}UpdateFunc
}

func {{.Name.Pascal}}StateFromString(s string) (*{{.Name.Pascal}}StateInstance, error) {
    state := &{{.Name.Pascal}}StateInstance{
        mu: &sync.RWMutex{},
        subs: map[uuid.UUID]{{.Name.Pascal}}UpdateFunc{},
    }

    switch s {
    {{range .States -}}
    case "{{.Pascal}}":
        state.current = {{$.Name.Pascal}}{{.Pascal}}
    {{end -}}
    default:
        return nil, fmt.Errorf("unknown state: %s", s)
    }

    return state, nil
}

func (s *{{.Name.Pascal}}StateInstance) String() string {
    switch s.current {
    {{range .States -}}
    case {{$.Name.Pascal}}{{.Pascal}}:
        return "{{.Pascal}}"
    {{end -}}
    default:
        panic(fmt.Sprintf("unknown state: %d", s.current))
    }
}

func (s *{{.Name.Pascal}}StateInstance) Current() {{.Name.Pascal}}State {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.current
}

func (s *{{.Name.Pascal}}StateInstance) Subscribe(f {{.Name.Pascal}}UpdateFunc) func() {
    s.mu.Lock()
    defer s.mu.Unlock()
    id := uuid.New()
    s.subs[id] = f
    return func() {
        s.mu.Lock()
        defer s.mu.Unlock()
        delete(s.subs, id)
    }
}

func (s *{{.Name.Pascal}}StateInstance) Set(state {{.Name.Pascal}}State) error {
    if state == s.current {
        return nil
    }

    s.mu.Lock()
    defer s.mu.Unlock()
    if s.isPaused {
        return fmt.Errorf("state machine is paused")
    }
    s.current = state
    for _, f := range s.subs {
        if err := f(state, s.isPaused); err != nil {
            return err
        }
    }
    return nil
}

func (s *{{.Name.Pascal}}StateInstance) Pause() {
    if s.isPaused {
        return
    }
    s.mu.Lock()
    defer s.mu.Unlock()
    s.isPaused = true
    for _, f := range s.subs {
        f(s.current, s.isPaused)
    }
}

func (s *{{.Name.Pascal}}StateInstance) Resume() {
    if !s.isPaused {
        return
    }
    s.mu.Lock()
    defer s.mu.Unlock()
    s.isPaused = false
    for _, f := range s.subs {
        f(s.current, s.isPaused)
    }
}
