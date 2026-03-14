package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"md2img/internal/cmd"
	"md2img/internal/controller"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	versionFlag   = flag.BoolP("version", "V", false, "查看当前版本后退出")
	installFlag   = flag.BoolP("install", "I", false, "安装服务并退出")
	uninstallFlag = flag.BoolP("uninstall", "U", false, "卸载服务并退出")
	hostFlag      = flag.String("host", "", "服务地址（支持多个地址逗号分隔，默认为所有地址）")
	portFlag      = flag.StringP("port", "p", "8080", "服务端口")

	logOut = log.New(os.Stdout, "", log.LstdFlags)
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	exit, err := doFlag()
	if err != nil {
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

	logOut.Println("POST /markdown?width=0&mobile=false&html=false&wait=1s")
	logOut.Println("POST /html?width=0&mobile=false&wait=1s")
	logOut.Println("POST /url?width=0&mobile=false&wait=1s")

	hosts := parseHosts(*hostFlag)
	errC := make(chan error, len(hosts))
	server := &http.Server{Handler: mux}

	for _, host := range hosts {
		addr := net.JoinHostPort(host, *portFlag)

		ln, err := net.Listen("tcp", addr)
		if err != nil {
			errC <- fmt.Errorf("listen %s: %w", addr, err)
			break
		}

		logOut.Println("http server started listening on", addr)

		// 同时监听多个 host，任一监听失败则退出。
		go func(listener net.Listener, listenAddr string) {
			if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errC <- fmt.Errorf("serve %s: %w", listenAddr, err)
			}
		}(ln, addr)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case err = <-errC:
	case sig := <-quit:
		log.Println("received signal:", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return errors.Join(err, server.Shutdown(ctx))
}

func parseHosts(rawHost string) []string {
	rawHost = strings.TrimSpace(rawHost)
	if rawHost == "" {
		return []string{""}
	}

	parts := strings.Split(rawHost, ",")
	hosts := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		host := strings.TrimSpace(part)
		if host == "" {
			continue
		}
		if _, ok := seen[host]; ok {
			continue
		}
		seen[host] = struct{}{}
		hosts = append(hosts, host)
	}

	if len(hosts) == 0 {
		return []string{""}
	}

	return hosts
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
