package serve

import (
	"fmt"
	"io"
	"merodi/src/core/lua"
	"merodi/src/core/state"
	"merodi/src/utils"
	"merodi/src/utils/log"
	"net/http"
)

var (
	handle = utils.Handle
	check  = utils.Check
)

type Serve struct {
	State *state.State
}

func (self *Serve) handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn(log.Title("HTTP"), "%s", err)
		return
	}
	requests := lua.Request{}
	requests.Method = r.Method
	requests.Body = string(body)
	requests.RemoteAddr = r.RemoteAddr
	response, err := self.State.Lua.RunRoute(r.URL.Path, requests)
	if err != nil {
		log.Warn(log.Title("HTTP"), "%s", err)
		return
	}
	w.Header().Set("Content-Type", response.Mime)
	w.WriteHeader(response.Status)
	_, err = w.Write([]byte(response.Body))
	if err != nil {
		log.Warn(log.Title("HTTP"), "%s", err)
	}
}
func (self *Serve) StartHTTP() (err error) {
	defer handle(&err)
	host := self.State.Lua.Http.Settings.Host
	port := self.State.Lua.Http.Settings.Port
	mux := http.NewServeMux()
	mux.HandleFunc("/", self.handler)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), mux)
}
func (self *Serve) Run(state *state.State, args *[]string) (err error) {
	defer handle(&err)
	self.State = state
	check(self.StartHTTP())
	return err
}
