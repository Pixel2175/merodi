package watch

import (
	"errors"
	"merodi/src/core/state"
	"merodi/src/parser"
	. "merodi/src/utils/errs"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	State *state.State
}

func (self *Watcher) Run(state *state.State, opt *parser.Option) (err error) {
	defer Handle(&err)
	self.State = state

	watcher := CheckV(fsnotify.NewWatcher())
	defer watcher.Close()

	if len(state.Watch.Add) == 0 {
		CheckE(watcher.Add(state.Config.Tree.Markdown))
	} else {
		for _, file := range self.State.Watch.Add {
			_, statErr := os.Stat(file)
			if os.IsNotExist(statErr) {
				CheckE(statErr)
			}
			CheckE(watcher.Add(file))
		}
	}

	CheckE(self.State.RunHook("on_start_watching"))

	lastEvents := make(map[string]time.Time)
	const debounce = 100 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			now := time.Now()
			if last, ok := lastEvents[event.Name]; ok && now.Sub(last) < debounce {
				continue
			}

			lastEvents[event.Name] = now
			CheckE(state.RunHook("on_file_changed", event.Op.String(), event.Name))

		case werr, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			if !errors.Is(werr, ErrAborted) {
				CheckE(werr)
			}
		}
	}
}
