
package cmd

import (
	"demo/pkg/logger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)
var name string
var age int


var (
	Hosts          string        // 全局hosts 变量
	RunICMP        bool          // 是否执行ICMP
	SurvivalHost   []string       //主机存活切片
	SurvHosts      int           //多少个主机存活
    redirect_Url   []string        //存活主机web url完整路径 切片
)
var rootCmd = &cobra.Command{
	Use:   "B1scan",
	Short: "A test B1scan",
	Long:  logger.Red(`                                                               
         __________________                                           
  |  ==c(______(o(______(_()  |                                             
  |             )=\           |                                          
  |            // \\          |                                           
  |           //   \\         |                                         
  |          //     \\        |                                         
  |         //       \\       |                                         
  |        //  B1scan \\                                                                           
                      .         `)+logger.Purple("～～"),
	Run: func(cmd *cobra.Command, args []string) {
		if len(name) == 0 {
			cmd.Help()
			fmt.Println(Hosts)
			return
		}

	},
}


func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&Hosts, "H", "H", "", "扫描的主机ip设置的-H参数格式    eg:192.168.2.10~~192.168.2.1-10~~192.168.2.1/24")
}



