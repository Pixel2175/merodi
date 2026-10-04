package watch

import (
	"errors"
	"merodi/src/core/state"
	. "merodi/src/utils/errs"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct{}

func (self *Watcher) Run(state *state.State, args *[]string) (err error) {
	defer Handle(&err)

	watcher := CheckV(fsnotify.NewWatcher())
	defer watcher.Close()

	if len(state.Lua.Watch.Add) == 0 {
		CheckE(watcher.Add(state.Config.Tree.Markdown))
	} else {
		for _, file := range state.Lua.Watch.Add {
			_, statErr := os.Stat(file)
			if os.IsNotExist(statErr) {
				CheckE(statErr)
			}
			CheckE(watcher.Add(file))
		}
	}

	CheckE(state.Lua.RunHook("on_start_watching"))

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
			CheckE(state.Lua.RunHook("on_file_changed", event.Op.String(), event.Name))

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
