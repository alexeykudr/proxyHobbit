package handler

import (
	"context"
	"fmt"
	"time"

	// "log"
	"os"
	"os/exec"
	// "time"
)

func ping(ctx context.Context, ip string) {
	cmd := exec.CommandContext(ctx, "ping", ip, "-t")
	cmd.Stdout = os.Stdin

	go cmd.Run()
	<-ctx.Done()
	fmt.Println("Pid of procces", cmd.Process.Pid)
	cmd.Process.Kill()
}

// /home/mac/goServer/reconect.sh

func (h *Handler) execProxyCommand() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "/bin/sh", "/home/mac/goServer/test.sh", "24").Run(); err != nil {
		fmt.Println("Error by context executor router rebooter!")
	}

}
