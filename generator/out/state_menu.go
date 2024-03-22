package out

import (
    "fmt"
    "sync"
    "github.com/google/uuid"
)

type MenuState int

const (
    MenuMain MenuState = 0
    MenuOptions MenuState = 1
    MenuQuit MenuState = 2
    )

type MenuUpdateFunc func(s MenuState, isPaused bool) error

type MenuStateInstance struct {
    current MenuState
    isPaused bool
    mu *sync.RWMutex
    subs map[uuid.UUID]MenuUpdateFunc
}

func MenuStateFromString(s string) (*MenuStateInstance, error) {
    state := &MenuStateInstance{
        mu: &sync.RWMutex{},
        subs: map[uuid.UUID]MenuUpdateFunc{},
    }

    switch s {
    case "Main":
        state.current = MenuMain
    case "Options":
        state.current = MenuOptions
    case "Quit":
        state.current = MenuQuit
    default:
        return nil, fmt.Errorf("unknown state: %s", s)
    }

    return state, nil
}

func (s *MenuStateInstance) String() string {
    switch s.current {
    case MenuMain:
        return "Main"
    case MenuOptions:
        return "Options"
    case MenuQuit:
        return "Quit"
    default:
        panic(fmt.Sprintf("unknown state: %d", s.current))
    }
}

func (s *MenuStateInstance) Current() MenuState {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.current
}

func (s *MenuStateInstance) Subscribe(f MenuUpdateFunc) func() {
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

func (s *MenuStateInstance) Set(state MenuState) error {
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

func (s *MenuStateInstance) Pause() {
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

func (s *MenuStateInstance) Resume() {
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
