// ServeHTTP implementation.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Add("Allow", "POST") // RFC 9110.
		core.WriteError(w, core.Error(http.StatusMethodNotAllowed, "only POST is allowed"))
		return
	}
	ctx := core.NewRequestContext(r.Context(), r)
	var res any
	var err error
	switch r.URL.Path {
	case "/add_item":
		var in AddItemInput
		var out AddItemOutput
		err = core.ReadRequest(r, &in)
		if err != nil {
			break
		}
		out, err = s.AddItem(ctx, in)
		res = out
	case "/get_items":
		var in GetItemsInput
		var out GetItemsOutput
		out, err = s.GetItems(ctx, in)
		res = out
	case "/remove_item":
		var in RemoveItemInput
		var out RemoveItemOutput
		err = core.ReadRequest(r, &in)
		if err != nil {
			break
		}
		out, err = s.RemoveItem(ctx, in)
		res = out
	default:
		err = core.Error(http.StatusNotFound, "method not found")
	}

	if err != nil {
		core.WriteError(w, err)
		return
	}

	core.WriteResponse(w, res)
}
