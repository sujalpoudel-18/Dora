package mysql

type Dialect struct{}

func (d *Dialect) Name() string { return "mysql" }
func (d *Dialect) DefaultPort() int { return 3306 }
