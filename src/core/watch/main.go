package watch

import (
	"merodi/src/core/state"
	"merodi/src/utils"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct{}

var (
	handle = utils.Handle
	check  = utils.Check
)

func (self *Watcher) Run(state *state.State, args *[]string) (err error) {
	defer handle(&err)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if len(state.Lua.Watch.Add) == 0 {
		check(watcher.Add(state.Config.Tree.Markdown))
	} else {
		for _, file := range state.Lua.Watch.Add {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return err
			} else {
				check(watcher.Add(file))
			}
		}
	}

	check(state.Lua.RunHook("on_start_watching"))

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
			check(state.Lua.RunHook("on_file_changed", event.Op.String(), event.Name))

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}

			if err != utils.ErrAborted {
				return err
			}
		}
	}
}
