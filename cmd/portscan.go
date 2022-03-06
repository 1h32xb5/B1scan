
package cmd

import (
	"fmt"
	"net"
	"sort"

	"github.com/spf13/cobra"
)
var(
 openports   []int
 closedports []int
)
var portscanCmd = &cobra.Command{
	Use:   "portscan",
	Short: "基于host主机的相应端口开放扫描",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("portscan called")
        fmt.Println(Hosts)
//		Ping()
		Portscan()
	},
}

func init() {
	rootCmd.AddCommand(portscanCmd)
	portscanCmd.Flags().StringVarP(&Hosts, "hosts", "H", "", "设置你要扫描端口 开放的主机")
}

func bcd(ports chan int,results chan int){        //channel类型、 wg指针    返回的一个结果channel
	for p:= range ports{                        //最后要关闭channel 不然一直在这里卡着
		address := fmt.Sprintf(Hosts+":"+"%d",p)
		conn,err := net.Dial("tcp",address)

		if err != nil {
			results <- 0
			continue                                   //return 退出执行 返回函数
		}
		results <- p
		conn.Close()
	}
}
func Portscan(){
	fmt.Printf("\033[1;31;40m%s\033[0m\n","Scan port open in progress...(可能存在漏扫 建议多扫几遍)")

	ports := make(chan int,100)              //开启100个channel缓冲通道
	results := make(chan int)            //没有缓冲 只有一个channel 因为main goroutine只有一个


	for i:= 0;i<cap(ports);i++{
		go bcd(ports,results)
	}

	go func() {                            //又一个goroutine 与下面的收集结果的分开 不会堵塞 下面收集结果的逻辑进行
		for i:=1;i<65535;i++{                   //分配工作
			ports <- i
		}

	}()
	for i:=1;i<100;i++{                //收集结果             收集结果的这个工作一定要在分配工作之前进行
		port := <- results
		if port != 0{
			openports = append(openports,port)
		}else{
			closedports = append(closedports,port)
		}
	}


	sort.Ints(openports)
	sort.Ints(closedports)

	for _,port := range openports{
		fmt.Printf(Hosts+":%d Open\n",port)
	}
}