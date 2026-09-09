package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func gluaToGo(v lua.LValue) any {
	switch v.Type() {
	case lua.LTString:
		return v.String()
	case lua.LTNumber:
		return float64(v.(lua.LNumber))
	case lua.LTBool:
		return bool(v.(lua.LBool))
	case lua.LTTable:
		return gluaTableToGo(v.(*lua.LTable))
	default:
		return v.String()
	}
}

func gluaTableToGo(t *lua.LTable) any {
	maxN := t.Len()
	if maxN > 0 {
		arr := make([]any, 0, maxN)
		for i := 1; i <= maxN; i++ {
			arr = append(arr, gluaToGo(t.RawGetInt(i)))
		}
		return arr
	}

	m := make(map[string]any)
	t.ForEach(func(k, v lua.LValue) {
		m[k.String()] = gluaToGo(v)
	})
	return m
}

func goToGlua(v any) lua.LValue {
	switch val := v.(type) {
	case string:
		return lua.LString(val)
	case float64:
		return lua.LNumber(val)
	case bool:
		return lua.LBool(val)
	case nil:
		return lua.LNil
	default:
		return lua.LNil
	}
}
