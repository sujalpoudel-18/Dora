package sqlite

type Dialect struct{}

func (d *Dialect) Name() string { return "sqlite" }
func (d *Dialect) DefaultPort() int { return 0 }
