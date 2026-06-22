type MapWGLog struct {
}

func (h *MapWGLog) Type() reflect.Type {
	return reflect.TypeOf((*gcn3.MapWGReq)(nil))
}

func (h *MapWGLog) Pos() core.HookPos {
	return core.AfterEvent
}

func (h *MapWGLog) Func(item interface{}, domain core.Hookable, info interface{}) {
	...	
}
