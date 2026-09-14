package mssql

type Dialect struct{}

func (d *Dialect) Name() string { return "mssql" }
func (d *Dialect) DefaultPort() int { return 1433 }
