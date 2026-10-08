package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"golang.org/x/term"
)

const refreshInterval = 3 * time.Second

// realCodexBin retorna o caminho absoluto do binário node do Codex,
// evitando loop com o próprio wrapper (que é instalado como "codex").
func realCodexBin() string {
	// Primeiro: CODEX_REAL_BIN no ambiente (set pelo install.sh).
	if v := os.Getenv("CODEX_REAL_BIN"); v != "" {
		return v
	}
	// Segundo: procura codex.js no path do mise/node.
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, ".local/share/mise/installs/node/22/bin/codex.js"),
		filepath.Join(home, ".local/share/mise/installs/node/22/lib/node_modules/@openai/codex/bin/codex.js"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	// Fallback: resolve "codex" ignorando alias do shell.
	if out, err := exec.LookPath("codex"); err == nil {
		return out
	}
	return "codex"
}

func nodebin() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, ".local/share/mise/installs/node/22/bin/node"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	if out, err := exec.LookPath("node"); err == nil {
		return out
	}
	return "node"
}

func main() {
	// Modo "stats" — imprime custo atual e sai (usado por scripts externos).
	if len(os.Args) > 1 && os.Args[1] == "stats" {
		sess := ReadSessionStats()
		today := ReadTodayStats()
		fmt.Println(RenderBar(sess, today))
		return
	}

	codexScript := realCodexBin()
	node := nodebin()

	// Monta args: node <codex.js> [args do usuário exceto o próprio wrapper]
	args := append([]string{codexScript}, os.Args[1:]...)

	// Passa --dangerously-bypass-approvals-and-sandbox se ainda não vier nos args.
	hasBypass := false
	for _, a := range os.Args[1:] {
		if strings.Contains(a, "dangerously") {
			hasBypass = true
			break
		}
	}
	if !hasBypass {
		args = append([]string{codexScript, "--dangerously-bypass-approvals-and-sandbox"}, os.Args[1:]...)
	}

	cmd := exec.Command(node, args...)
	cmd.Env = append(os.Environ(), "CODEX_REAL_BIN="+codexScript)

	// Cria PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-statusline: pty.Start: %v\n", err)
		os.Exit(1)
	}
	defer ptmx.Close()

	// Propaga redimensionamento de janela.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if sz, err := pty.GetsizeFull(os.Stdin); err == nil {
				_ = pty.Setsize(ptmx, sz)
			}
		}
	}()
	ch <- syscall.SIGWINCH // dispara uma vez pra sincronizar tamanho inicial

	// Coloca stdin em raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-statusline: MakeRaw: %v\n", err)
	} else {
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	// Stdin → PTY.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()

	// PTY → Stdout (output do Codex).
	go func() { _, _ = io.Copy(os.Stdout, ptmx) }()

	// Goroutine de atualização da barra de status.
	go func() {
		// Aguarda o TUI inicializar antes de exibir a barra.
		time.Sleep(2 * time.Second)
		for {
			sess := ReadSessionStats()
			today := ReadTodayStats()
			bar := PrintBar(RenderBar(sess, today))
			_, _ = os.Stdout.WriteString(bar)
			time.Sleep(refreshInterval)
		}
	}()

	_ = cmd.Wait()

	// Restaura terminal antes de sair.
	if oldState != nil {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}
	// Limpa a linha de status.
	fmt.Print("\r\x1b[K")

	if exitErr, ok := cmd.ProcessState.Sys().(interface{ ExitCode() int }); ok {
		os.Exit(exitErr.ExitCode())
	}
}
