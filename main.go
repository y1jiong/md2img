package main

import (
	flag "github.com/spf13/pflag"
	"log"
	"md2img/internal/cmd"
	"md2img/internal/controller"
	"net/http"
)

var (
	versionFlag   = flag.BoolP("version", "V", false, "查看当前版本后退出")
	installFlag   = flag.BoolP("install", "I", false, "安装服务并退出")
	uninstallFlag = flag.BoolP("uninstall", "U", false, "卸载服务并退出")
	addressFlag   = flag.StringP("address", "a", ":8080", "服务地址")
)

func main() {
	exit, err := doFlag()
	if err != nil {
		log.Fatal(err)
		return
	}
	if exit {
		return
	}

	mux := http.NewServeMux()

	// 设置路由
	mux.HandleFunc("POST /markdown", controller.Markdown)
	mux.HandleFunc("POST /html", controller.HTML)
	mux.HandleFunc("POST /url", controller.URL)

	log.Println("http server started listening on", *addressFlag)
	log.Println("POST /markdown?width=0&mobile=false&html=false")
	log.Println("POST /html?width=0&mobile=false")
	log.Println("POST /url?width=0&mobile=false")

	// 启动
	log.Fatal(http.ListenAndServe(*addressFlag, mux))
}

func doFlag() (exit bool, err error) {
	flag.Parse()
	if *versionFlag {
		cmd.PrintVersion()
		return true, nil
	}
	if *installFlag {
		return true, cmd.Install()
	}
	if *uninstallFlag {
		return true, cmd.Uninstall()
	}
	return
}
