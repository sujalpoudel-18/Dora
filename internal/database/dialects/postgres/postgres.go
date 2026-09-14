package postgres

type Dialect struct{}

func (d *Dialect) Name() string { return "postgresql" }
func (d *Dialect) DefaultPort() int { return 5432 }
