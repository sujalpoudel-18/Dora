package oracle

type Dialect struct{}

func (d *Dialect) Name() string { return "oracle" }
func (d *Dialect) DefaultPort() int { return 1521 }
