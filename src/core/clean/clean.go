package clean

import (
	"merodi/src/config"
	"merodi/src/core/state"
	. "merodi/src/utils/errs"
	"os"
)

type Clean struct {
	State *state.State
}

func (self *Clean) Run(state *state.State, args *[]string) (err error) {
	defer Handle(&err)
	self.State = state

	if self.State.Build.Mode == config.Release {
		CheckE(os.RemoveAll(self.State.Config.Tree.ReleaseDest))
	} else {
		CheckE(os.RemoveAll(self.State.Config.Tree.DraftDest))
	}
	return
}
