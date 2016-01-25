package proj

var defs map[string]*Proj

func addDef(name, def string) error {
	proj, err := parse(def)
	if err != nil {
		return err
	}
	defs[name] = proj
	return nil
}
