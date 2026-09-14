package redis

type Dialect struct{}

func (d *Dialect) Name() string { return "redis" }
func (d *Dialect) DefaultPort() int { return 6379 }
