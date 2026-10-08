# PTY Wrapper — codex-statusline

## Objetivo
Injetar linha de custo/tokens no TUI do Codex sem depender dos dados do gateway Superlogica.

## Arquivos
- `sessions.go` — leitura de `~/.codex/sessions/*.jsonl`, cálculo de custo
- `statusline.go` — renderização ANSI da linha de status
- `main.go` — PTY wrapper: lança `codex`, redireciona I/O, injeta status no rodapé

## Approach
PTY wrapper em Go usando `github.com/creack/pty`:
1. `main.go` lança o binário real do Codex num PTY
2. Copia stdin/stdout/stderr transparentemente
3. A cada turno (detectado por silêncio no output), lê `~/.codex/sessions/` e 
   atualiza a linha de status na última linha do terminal via ANSI escape codes

## Linha de status (formato)
```
 ~/projeto · main · gpt-6-luna · $0.05 sessão · $0.12 hoje · 17.4K tokens
```

## Instalação
`install.sh` compila e substitui o alias no `.zshrc`
