package cmd

import (
	"database/sql"
	"demo/cmd/config"
	"demo/pkg/logger"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/spf13/cobra"
	"log"
)
var mssqlCmd = &cobra.Command{
	Use:   "mssql",
	Short: "对内网mssql数据库进行弱口令检测",
	Run: func(cmd *cobra.Command, args []string) {
		Choice()
		passwords := config.Passwords
		users:=config.Userdict["mssql"]
		for _, user := range users {
			for _, password := range passwords {
				for _, ip := range SurvivalHost {
					success := Burtemssql(user,password,ip)
					if success == true {
						log.Println(logger.LightGreen(ip+" "+user+" "+password+" "+logger.LightGreen(success)))
					}
					if success {
						c:=fmt.Sprintf("破解%v成功，用户名是%v,密码是%v\n", ip, user, password)
						fmt.Println(logger.Purple(c))
					}
				}
			}
		}
	},
}

func Burtemssql(user string,password string,ip string) bool{
	var DB *sql.DB
	connString := fmt.Sprintf("sqlserver://%v:%v@%v:%v/?connection&timeout=%v&encrypt=disable", user, password,ip, 1433, 3)
	DB, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	DB.SetConnMaxLifetime(10)          //设置最大连接数
	DB.SetMaxIdleConns(10)            //设置上数据库最大闲置连接数
	if err := DB.Ping(); err != nil{           //连接数据库 验证连接
		return false
	}else {
		fmt.Println("打开数据库成功")
		return true
	}
}
func init() {
	rootCmd.AddCommand(mssqlCmd)
	mssqlCmd.Flags().StringVarP(&Hosts, "hosts", "H", "", "设置你要扫描mysqld爆破的主机")
}

