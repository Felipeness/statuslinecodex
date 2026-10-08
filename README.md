# Codex Statusline Studio

Editor visual offline para montar a statusline nativa do Codex. Inclui um wrapper Go que injeta uma segunda linha de custo calculado localmente — útil quando o gateway não envia dados de limite (ex: Superlogica LLM Gateway).

## Wrapper de custo (`codex-statusline`)

O Codex não expõe custo quando usado via gateway corporativo. O wrapper lê os arquivos `~/.codex/sessions/*.jsonl` e injeta uma linha extra no rodapé do TUI:

```
▸ $0.055 sessão · $0.088 hoje · 14.5K tokens · 96% cache · 1 turno
```

### Instalação

**Pré-requisito:** Go instalado (`mise install go` ou `brew install go`)

```bash
git clone https://github.com/Felipeness/statuslinecodex ~/statuslinecodex
cd ~/statuslinecodex
bash install.sh
source ~/.zshrc
```

O `install.sh` compila o binário e atualiza o alias `codex` no `.zshrc` para usar o wrapper automaticamente.

### Uso

Use `codex` normalmente — o wrapper é transparente. Para ver o custo atual sem abrir o Codex:

```bash
codex-statusline stats
```

### Como funciona

O wrapper lança o Codex num PTY e redireciona I/O de forma transparente. A cada 3 segundos lê os `token_usage_record` dos arquivos de sessão e recalcula o custo usando os preços públicos do modelo:

| Tokens | Preço por milhão |
|--------|-----------------|
| Input  | $3,00 |
| Output | $15,00 |
| Cache read | $0,75 |
| Cache write | $3,75 |
| Reasoning | $15,00 |

> No Plus/Pro ou via gateway corporativo o valor é o que as chamadas custariam via API direta, não uma cobrança real.

---

## Studio visual

Abra [`index.html`](index.html) em um navegador. Não precisa de servidor, dependências ou acesso à internet.

O arquivo [`config.toml`](config.toml) traz a configuração padrão usada pelo editor. Para aplicá-la, incorpore a tabela `[tui]` ao seu `~/.codex/config.toml` ou ao caminho definido por `CODEX_HOME`, preservando as outras opções que já existirem. Reinicie o Codex para carregar a alteração.

## Itens nativos de custo

O Codex não aceita renderer customizado em `tui.status_line`, só identificadores nativos. Para custo, existem três caminhos:

- `estimated-thread-cost` e `thread-credits`: custo e créditos estimados da thread, **só em workspace Enterprise**. Fora disso o item fica oculto.
- `five-hour-limit` e `weekly-limit`: quanto resta da cota no Plus/Pro. Nesses planos é a métrica de gasto que importa, já que não há fatura por token.
- `used-tokens`, `total-input-tokens` e `total-output-tokens`: volume de tokens da sessão.

### Alternativa: ccusage

O [`ccusage`](https://github.com/ryoppippi/ccusage) lê os arquivos de sessão e calcula custo via linha de comando:

```bash
npx ccusage@latest codex daily     # por dia
npx ccusage@latest codex monthly   # por mês
npx ccusage@latest codex session   # por sessão
```

## Itens disponíveis

O Studio lista todos os itens de `StatusLineItem` do TUI (`codex-rs/tui/src/bottom_pane/status_line_setup.rs`): diretório, projeto, branch, PR, mudanças da branch, modelo, raciocínio, fast mode, contexto (restante, usado, janela), tokens, limites de 5h e semanal, custo e créditos Enterprise, estado, progresso do plano, permissões, aprovação, thread, hostname, versão e afins. Itens sem dado ficam ocultos.

## Configuração padrão

```toml
[tui]
status_line = ["current-dir", "git-branch", "model-with-reasoning", "context-remaining", "used-tokens", "five-hour-limit", "weekly-limit", "estimated-thread-cost"]
terminal_title = ["spinner", "project"]
```

Confira a [referência oficial de configuração do Codex](https://developers.openai.com/codex/config-file/config-reference.md) para detalhes sobre `tui.status_line`.
