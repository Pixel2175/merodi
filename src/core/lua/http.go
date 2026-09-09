package lua

import (
	"merodi/src/utils/log"

	glua "github.com/yuin/gopher-lua"
)

type Request struct {
	RemoteAddr string
	Body       string
	Method     string
}

type Response struct {
	Mime   string
	Body   string
	Status int
}

type Settings struct {
	Host string
	Port string
}

type Http struct {
	Settings Settings
	Table    *glua.LTable
	Route    map[string]*glua.LFunction
}

func (self *Lua) registerHttp() {
	self.Http.Table = self.Context.NewTable()

	self.Http.Route = make(map[string]*glua.LFunction)
	self.httpSettings()
	self.httpRoute()
	self.Merodi.RawSetString("http", self.Http.Table)
}

func (self *Lua) httpSettings() {
	settingsTable := self.Context.NewTable()

	settingsTable.RawSetString("host", self.makeGetSetTable(
		func() string { return self.Http.Settings.Host },
		func(v string) { self.Http.Settings.Host = v },
	))

	settingsTable.RawSetString("port", self.makeGetSetTable(
		func() string { return self.Http.Settings.Port },
		func(v string) { self.Http.Settings.Port = v },
	))

	self.Http.Table.RawSetString("settings", settingsTable)
}

func (self *Lua) httpRequest(req Request) {
	requestTable := self.Context.NewTable()

	self.Context.SetFuncs(requestTable, map[string]glua.LGFunction{
		"body": func(L *glua.LState) int {
			L.Push(glua.LString(req.Body))
			return 1
		},
		"remote_addr": func(L *glua.LState) int {
			L.Push(glua.LString(req.RemoteAddr))
			return 1
		},
		"method": func(L *glua.LState) int {
			L.Push(glua.LString(req.Method))
			return 1
		},
	})

	self.Http.Table.RawSetString("request", requestTable)
}

func (self *Lua) httpRoute() {
	self.Http.Table.RawSetString("route", self.Context.NewFunction(func(L *glua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)
		self.Http.Route[path] = fn
		return 0
	}))
}

func (self *Lua) RunRoute(path string, req Request) (Response, error) {
	fn, ok := self.Http.Route[path]
	if !ok {
		return Response{Mime: "text/plain", Body: "this page does not exists", Status: 404}, nil
	}

	self.httpRequest(req)
	err := self.Context.CallByParam(glua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	})

	if err != nil {
		log.Warn(log.Title("HTTP"), "route: %s", err)
		return Response{}, err
	}

	ret := self.Context.Get(-1)
	self.Context.Pop(1)

	resp := Response{
		Mime:   "text/plain",
		Body:   "",
		Status: 200,
	}

	if t, ok := ret.(*glua.LTable); ok {
		if v := t.RawGetString("mime"); v != glua.LNil {
			resp.Mime = v.String()
		}
		if v := t.RawGetString("body"); v != glua.LNil {
			resp.Body = v.String()
		}
		if v := t.RawGetString("status"); v != glua.LNil {
			if n, ok := v.(glua.LNumber); ok {
				resp.Status = int(n)
			}
		}
	}

	return resp, nil
}
