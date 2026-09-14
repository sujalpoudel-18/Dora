package mongodb

type Dialect struct{}

func (d *Dialect) Name() string { return "mongodb" }
func (d *Dialect) DefaultPort() int { return 27017 }
