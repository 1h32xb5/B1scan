package cmd

import (
	"database/sql"
	"demo/cmd/config"
	"demo/pkg/logger"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

const (
	port = "3306"
	dbName = "information_schema"            //数据库 名字  下面打开这个数据库     information_schema
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "对内网mysql数据库进行弱口令检测",
	Run: func(cmd *cobra.Command, args []string) {
		Choice()
		passwords := config.Passwords
		users:=config.Userdict["mysql"]
		for _, user := range users {                              //先是以password遍历ip，然后以user遍历password
			for _, password := range passwords {
				for _, ip := range SurvivalHost {
					success := Burtemysql(user,password,ip)
					if success == true {                           //判断是否是成功 如果成功则 高亮 true
						log.Println(logger.LightGreen(ip+" "+user+" "+password+" "+logger.LightGreen(success)))
					}
					if success {
						c:=fmt.Sprintf("破解%v成功，用户名是%v,密码是%v\n", ip, user, password)
						fmt.Println(logger.LightBlue(c))
					}
				}
			}
		}
	},
}
func Burtemysql(user string,password string,ip string) bool{
		var DB *sql.DB
		path:=strings.Join([]string{user, ":", password, "@tcp(",ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
		DB, _ = sql.Open("mysql", path)
		DB.SetConnMaxLifetime(5)
		DB.SetMaxIdleConns(5)
		if err := DB.Ping(); err != nil{           //连接数据库 验证连接
			return false
		}else {
			//rows,_:=DB.Query("SELECT * FROM register ")       //获取所有数据
			////sql.NullString  如果是空字符串会报错
			////var id,name,a,b sql.NullString   //
			////fmt.Println(id.String,"--",name.String,a.String,b.String)
			//var id int
			//var usernames,passwords string
			//for rows.Next(){        //循环显示所有的数据
			//	rows.Scan(&id,&usernames,&passwords)     //
			//	fmt.Println(id,usernames,passwords)
			//}
			return true
		}
}
func init() {
	rootCmd.AddCommand(mysqlCmd)
	mysqlCmd.Flags().StringVarP(&Hosts, "hosts", "H", "", "设置你要扫描mysqld爆破的主机")
}


