package main

import (
	"fmt"
	"math"
	"strings"
)

// Cores ANSI — mesmo estilo do claude-statusline (256-color).
const (
	reset     = "\x1b[0m"
	bold      = "\x1b[1m"
	dim       = "\x1b[2m"
	fgGreen   = "\x1b[32m"
	fgYellow  = "\x1b[33m"
	fgCyan = "\x1b[36m"
	fgOrange  = "\x1b[38;5;208m"
	fgGray    = "\x1b[38;5;244m"
)

// saveCursor / restoreCursor / eraseToEOL são seqüências ANSI portáveis.
const (
	saveCursor    = "\x1b7"
	restoreCursor = "\x1b8"
	eraseToEOL    = "\x1b[K"
	moveToLastRow = "\x1b[999;1H" // cursor para linha muito além do fim → última linha
)

func formatTokens(n int) string {
	if n == 0 {
		return "0"
	}
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(n)/1e3)
	}
	return fmt.Sprintf("%d", n)
}

func costColor(usd float64) string {
	switch {
	case usd < 0.05:
		return fgGreen
	case usd < 0.20:
		return fgYellow
	case usd < 0.50:
		return fgOrange
	default:
		return "\x1b[31m" // vermelho
	}
}

func sep() string { return fgGray + " · " + reset }

// RenderBar monta a linha de status com os dados de sessão e dia.
// Formato: $0.03 sessão · $0.08 hoje · 17.4K tokens · 3 turnos
func RenderBar(session, today SessionStats) string {
	var b strings.Builder

	// custo sessão
	sessUSD := math.Round(session.CostUSD*1000) / 1000
	b.WriteString(costColor(sessUSD))
	b.WriteString(bold)
	fmt.Fprintf(&b, "$%.3f", sessUSD)
	b.WriteString(reset)
	b.WriteString(fgGray + " sessão" + reset)

	// custo hoje
	todayUSD := math.Round(today.CostUSD*1000) / 1000
	b.WriteString(sep())
	b.WriteString(costColor(todayUSD))
	fmt.Fprintf(&b, "$%.3f", todayUSD)
	b.WriteString(reset)
	b.WriteString(fgGray + " hoje" + reset)

	// tokens da sessão
	totalTok := session.InputTokens + session.OutputTokens
	if totalTok > 0 {
		b.WriteString(sep())
		fmt.Fprintf(&b, "%s%s%s", fgCyan, formatTokens(totalTok), reset)
		b.WriteString(fgGray + " tokens" + reset)
	}

	// cache hit rate (util pra saber se o contexto ta sendo reutilizado)
	if session.InputTokens > 0 {
		hitRate := float64(session.CachedRead) / float64(session.InputTokens) * 100
		if hitRate > 5 {
			b.WriteString(sep())
			fmt.Fprintf(&b, "%s%.0f%%%s", fgGreen, hitRate, reset)
			b.WriteString(fgGray + " cache" + reset)
		}
	}

	// turnos
	if session.Turns > 0 {
		b.WriteString(sep())
		turno := "turnos"
		if session.Turns == 1 {
			turno = "turno"
		}
		fmt.Fprintf(&b, "%s%d %s%s", dim, session.Turns, turno, reset)
	}

	return b.String()
}

// PrintBar escreve a barra na última linha do terminal e volta o cursor.
// Usa save/restore cursor pra não interferir com o TUI do Codex.
func PrintBar(bar string) string {
	return saveCursor +
		moveToLastRow +
		"\r" + eraseToEOL +
		fgGray + "▸ " + reset + bar +
		restoreCursor
}
