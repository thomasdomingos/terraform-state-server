package states

import (
  "crypto/sha256"
  "encoding/hex"
  "errors"
  "log"
  "strings"
)

type State struct {
  Name string
  Previous string
  Checksum string
  content string
}

func (s *State) computeChecksum() error {
  if s.Name == "" || s.content == "" {
    return errors.New("incomplete state, cannot compute checksum")
  }
  if s.Checksum != "" {
    return nil
  }
  h := sha256.New()
  h.Write([]byte(s.Name))
  h.Write([]byte(s.Previous))
  h.Write([]byte(s.content))
  s.Checksum = hex.EncodeToString(h.Sum(nil))
  return nil
}

func (s State) Verify() bool {
  h := sha256.New()
  h.Write([]byte(s.Name))
  h.Write([]byte(s.Previous))
  h.Write([]byte(s.content))
  return hex.EncodeToString(h.Sum(nil)) == s.Checksum
}

func NewState(name, content string) *State {
  state := new(State)
  state.Name = name
  state.Previous = strings.Repeat("0", sha256.BlockSize)
  state.content = content
  if err := state.computeChecksum(); err != nil {
    log.Println("ERROR", err)
    return nil
  }
  return state
}

func NextState(state State, newcontent string) *State {
  newstate := new(State)
  newstate.Name = state.Name
  newstate.Previous = state.Checksum
  newstate.content = newcontent
  if err := newstate.computeChecksum(); err != nil {
    log.Println("ERROR", err)
    return nil
  }
  return newstate
}

