package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gucooing/discordbot-go/pgk"
)

func main() {
	fmt.Println("    ___     __    _____    __")
	fmt.Println("   /   |   / /   / ___/   / /")
	fmt.Println("  / /| |  / /    \\__ \\   / /")
	fmt.Println(" / ___ | / /___ ___/ /  / /___")
	fmt.Println("/_/  |_|/_____//____/  /_____/")
	// 启动读取配置
	err := pgk.LoadConfig()
	if err != nil {
		if err == pgk.FileNotExist {
			p, _ := json.MarshalIndent(pgk.DefaultConfig, "", "  ")
			fmt.Printf("找不到配置文件，这是默认配置:\n%s\n", p)
			fmt.Printf("\n您可以将其保存到名为“config.json”的文件中并再次运行该程序\n")
			fmt.Printf("按 'Enter' 键退出 ...\n")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		} else {
			panic(err)
		}
	}
	go func() {
		for {
			pgk.DiscordBot()
			fmt.Printf("discord bot 失去连接 重连中 ...\n")
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		for {
			pgk.Wsls()
			fmt.Printf("Ws 服务器掉线了 重新开启中 ...\n")
			time.Sleep(5 * time.Second)
		}
	}()
	for {
		time.Sleep(50 * time.Second)
	}
}
