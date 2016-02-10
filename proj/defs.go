package proj

var defs map[string]*Proj

func addDef(name, def string) error {
	if defs == nil {
		defs = make(map[string]*Proj)
	}
	proj, err := Parse(def)
	if err != nil {
		return err
	}
	defs[name] = proj
	return nil
}
