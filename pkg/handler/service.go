package handler

import (
	"context"
	"fmt"
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

