package state

import (
	"github.com/hashicorp/nomad/client/structs"
	"github.com/hashicorp/nomad/helper"
	"github.com/hashicorp/nomad/plugins/drivers"
)

// LocalState is Task state which is persisted for use when restarting Nomad
// agents.
type LocalState struct {
	Hooks map[string]*HookState

	// DriverNetwork is the network information returned by the task
	// driver's Start method
	DriverNetwork *structs.DriverNetwork

	// TaskHandle is the handle used to reattach to the task during recovery
	TaskHandle *drivers.TaskHandle
}

func NewLocalState() *LocalState {
	return &LocalState{
		Hooks: make(map[string]*HookState),
	}
}

// Copy should be called with the lock held
func (s *LocalState) Copy() *LocalState {
	// Create a copy
	c := &LocalState{
		Hooks:         make(map[string]*HookState, len(s.Hooks)),
		DriverNetwork: s.DriverNetwork.Copy(),
		TaskHandle:    s.TaskHandle.Copy(),
	}

	// Copy the hooks
	for h, state := range s.Hooks {
		c.Hooks[h] = state.Copy()
	}

	return c
}

type HookState struct {
	// Prestart is true if the hook has run Prestart successfully and does
	// not need to run again
	PrestartDone bool
	Data         map[string]string
}

func (h *HookState) Copy() *HookState {
	c := new(HookState)
	*c = *h
	c.Data = helper.CopyMapStringString(c.Data)
	return c
}

func (h *HookState) Equal(o *HookState) bool {
	if h == nil || o == nil {
		return h == o
	}

	if h.PrestartDone != o.PrestartDone {
		return false
	}

	return helper.CompareMapStringString(h.Data, o.Data)
}
