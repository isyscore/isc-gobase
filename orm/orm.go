package orm

type DatasourceConfig struct {
	Username   string
	Password   string
	Host       string
	Port       int
	DriverName string
	DbName     string
	SuffixUrl  string
}
