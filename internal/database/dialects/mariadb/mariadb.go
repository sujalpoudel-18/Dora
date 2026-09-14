package mariadb

type Dialect struct{}

func (d *Dialect) Name() string { return "mariadb" }
func (d *Dialect) DefaultPort() int { return 3306 }
