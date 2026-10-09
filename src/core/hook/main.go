package hook

import (
	"merodi/src/core/state"
	"merodi/src/parser"
	. "merodi/src/utils/errs"
)

type Hook struct {
	State *state.State
}

func (self *Hook) Run(state *state.State, opt *parser.Option) (err error) {
	defer Handle(&err)
	self.State = state

	CheckE(self.State.RunHook(opt.Hook, opt.Args...))
	return
}
