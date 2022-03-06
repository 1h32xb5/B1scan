package config

import (
	"time"
)
type HostIn struct {
	Host      string
	Port      int
	Domain    string
	TimeOut   time.Duration
	PublicKey string
}

// 爆破的默认用户名
var Userdict = map[string][]string{
	"ftp":      {"kali", "ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"},
	"mysql":    {"root", "mysql"},
	"mssql":    {"sa", "sql"},
	"smb":      {"administrator", "admin", "guest"},
	"rdp":      {"administrator", "admin", "guest", "Oadmin"},
	"postgres": {"postgres", "admin"},
	"ssh":      {"root", "admin", "kali"},
	"mongodb":  {"root", "admin"},
	"redis":    {"root"},
}

// 爆破的默认密码

var Passwords = []string{"123456", "admin", "admin123", "root", "12312", "pass123", "pass@123", "930517", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#",  "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456", "1q2w3e", "kali"}


