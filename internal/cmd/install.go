package cmd

import (
	"errors"
	"log"
	"md2img/internal/consts"
	"os"
	"runtime"
)

const (
	installPath = "/etc/systemd/system/" + consts.ProjName + ".service"
)

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func Install() (err error) {
	if isWindows() {
		return errors.New("windows 暂不支持安装到系统")
	}
	// 注册系统服务
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	serviceContent := []byte(
		"[Unit]\n" +
			"Description=" + consts.ProjName + " Service\n" +
			"After=network-online.target\n\n" +
			"[Service]\n" +
			"Type=simple\n" +
			"WorkingDirectory=" + wd +
			"\nExecStart=" + wd + "/" + consts.ProjName +
			"\nRestart=on-failure\n" +
			"RestartSec=2\n\n" +
			"[Install]\n" +
			"WantedBy=multi-user.target\n",
	)
	if err = os.WriteFile(installPath, serviceContent, 0600); err != nil {
		return
	}
	log.Println("安装服务成功\n可以使用 systemctl 管理", consts.ProjName, "服务了")
	return
}

func Uninstall() (err error) {
	if isWindows() {
		return errors.New("windows 暂不支持安装到系统")
	}
	if err = os.Remove(installPath); err != nil {
		return
	}
	log.Println("卸载服务成功")
	return
}
