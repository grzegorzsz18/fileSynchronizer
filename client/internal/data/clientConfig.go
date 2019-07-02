package data

type ClientConfig struct {
	UserName          string
	UserPasswordHash  string
	ServerHost        string
	ServerPortRest    string
	ServerPortTCP     string
	RefreshFilesTime  int
	UserEncriptionKey string
}
