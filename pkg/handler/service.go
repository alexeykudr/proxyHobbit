package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func ping(ctx context.Context, ip string) {
	cmd := exec.CommandContext(ctx, "ping", ip, "-t")
	cmd.Stdout = os.Stdin

	go cmd.Run()
	<-ctx.Done()
	fmt.Println("Pid of procces", cmd.Process.Pid)
	cmd.Process.Kill()
}

func (h *Handler) execProxyCommand() {
	ctx, cancel := context.WithTimeout(context.Background(), 5500*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "sleep", "5").Run(); err != nil {
		log.Fatalf("Error in execProxyCommand %s", err.Error())
		fmt.Println("err!")
	}
}