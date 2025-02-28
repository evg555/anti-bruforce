package storage

const (
	Whitelist = "whitelist"
	Blacklist = "blacklist"
)

type Subnet struct {
	ID      int    `db:"id"`
	Address string `db:"subnet"`
}
