package cmd

import (
	"bytes"
	"crypto/tls"
	"demo/pkg/logger"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var resp_title string
var response_body string
var body_bytes string
var redirect1_Url string
var code int
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "一个全方位信息搜集的命令(主机存活/端口开放/webtitle信息)",
	Run: func(cmd *cobra.Command, args []string) {
		infoPing()                  //将存活主机放入到一个切片中
		start := time.Now()
		SurvHosts = len(SurvivalHost)
        fmt.Println("--------主机存活探测完成--------"+logger.Red("[除了不走icmp协议的主机]"))
		var wg sync.WaitGroup

		stringchan := make(chan string,200)


		for i:=0;i<len(SurvivalHost);i++{
			wg.Add(1)
			stringchan <- SurvivalHost[i]
		}
		for i:=0;i< 500;i++{
			go func ( ){
				a :=<- stringchan
				defer wg.Done()
				conn, _ := net.Dial("tcp", a)
				http_url := a + ":80"
				https_url :=a + ":443"
				conn, _ = net.Dial("tcp", http_url)
				conn1, _ := net.Dial("tcp", https_url)
				if conn != nil {
					redirect1_Url = "http://" + a
					redirect_Url = append(redirect_Url,redirect1_Url)
				} else if conn1 != nil {
					redirect1_Url = "https://" + a
					redirect_Url = append(redirect_Url,redirect1_Url)
				}
			}()
		}
		wg.Wait()
		tail()
		elapsed := time.Since(start)
		fmt.Println(elapsed)
	},
}
func tail(){
	for a,_:=range redirect_Url{
		request(a,redirect_Url,code,5)
		l:=fmt.Sprintf("[%d]",request(a,redirect_Url,code,3))
		fmt.Println(logger.Purple(l)+redirect_Url[a]+" [title:" + logger.Blue(resp_title)+"]")
	}
}
func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&Hosts, "hosts", "H", "", "设置你要全方位扫描的主机")     //起到的作用就是声明一个属于父命令的一个参数  但是不起实质作用
}
func request(i int,redirect_Url []string,code int,timeout int) int{
		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport {
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse                //重定向
			},
		}
		request, _ := http.NewRequest("GET", redirect_Url[i], nil)
		resp, _ := client.Do(request)
		body_bytes, _ := ioutil.ReadAll(resp.Body)
		response_body = string(body_bytes)
		grep_title := regexp.MustCompile("<title>(.*)</title>")
		if len(grep_title.FindStringSubmatch(response_body)) != 0 {
			resp_title = grep_title.FindStringSubmatch(response_body)[1]
		} else {
			resp_title = "None"
		}
		code = resp.StatusCode
	return code
}

func infoPing(){
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
	for i := 0; i < cap(cha); i++ {
		go Noutput(b, cha, &wg)
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
			cha <- a4
	}
	wg.Wait()
	close(cha)
}

func Noutput(b string, cha chan int, wg *sync.WaitGroup) {
	sysType := runtime.GOOS

	if sysType == "windows" {
		for p := range cha {
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




















