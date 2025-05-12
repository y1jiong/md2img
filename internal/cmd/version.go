package cmd

import (
	"fmt"
	"md2img/internal/consts"
)

func PrintVersion() {
	fmt.Println(consts.Description)
}
