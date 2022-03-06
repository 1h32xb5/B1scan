package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)



var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "icmp of ping hosts",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		Ping()
		elapsed := time.Since(start)
		fmt.Println(elapsed)
		defer func() {fmt.Println("c段扫描完成！")}()
	},
}
func init() {
		rootCmd.AddCommand(pingCmd)
		pingCmd.Flags().StringVarP(&Hosts, "hosts", "H", "", "设置你要扫描的主机")
}
func Ping() {
	flag.Parse()
	cha := make(chan int, 512)
	var wg sync.WaitGroup

	ip:=Hosts
	if ip == "" {
		fmt.Println("Please set a -H Parameter")
		fmt.Println("You may be need help (-h/-help")
		os.Exit(0)
	}
	fmt.Printf("\033[1;31;40m%s\033[0m\n","ICMP host survival scan in progress...")
	a := Hosts
	s := strings.Split(a, ".")
	b := s[0] + "." + s[1] + "." + s[2] + "."

	for i := 0; i < 512; i++ {
		go B(b, cha, &wg)
	}

	if strings.Contains(a, "/24"){
		for i := 1; i <= 255; i++ {
			wg.Add(1)
			cha <- i
		}
	}else {
		v := strings.Split(s[3],"-")
		j1, _ := strconv.Atoi(v[0])
		j2, _ := strconv.Atoi(v[1])
		for j := j1; j <= j2; j++{
			wg.Add(1)
			cha <- j
		}
	}
	wg.Wait()
	close(cha)
}

func B(b string, cha chan int, wg *sync.WaitGroup) {
	sysType := runtime.GOOS

	if sysType == "windows" {
		for p := range cha {
			//p 1-255 数字
			address := fmt.Sprintf("%s%d", b, p)

			cmd := exec.Command("cmd", "/c", "ping -n 1 "+address)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Run()

			if strings.Contains(out.String(), "TTL=") {
				fmt.Printf("%s主机存活\n", address)
				SurvivalHost=append(SurvivalHost,address )
			}
			wg.Done()
		}
	} else if sysType == "linux" {
		for p := range cha {
			address := fmt.Sprintf("%s%d", b, p)

			cmd := exec.Command("/bin/bash", "/c", "ping -n 1 "+address)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Run()

			if strings.Contains(out.String(), "TTL=") {
				fmt.Printf("%s主机存活\n", address)
				SurvivalHost=append(SurvivalHost,address )
			}
			wg.Done()
		}
	} else if sysType == "darwin"{
		for p:= range cha{
			address := fmt.Sprintf("%s%d",b,p)

			cmd:=exec.Command("/bin/bash", "-c", "ping -c 1 "+address)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Run()
			if strings.Contains(out.String(), "ttl=") {
				fmt.Printf("%s主机存活\n", address)
				SurvivalHost=append(SurvivalHost,address )         //把扫描的存活的主机 添加到全局切片里面
			}
			wg.Done()
		}
	}
}