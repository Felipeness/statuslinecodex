# Codex Statusline Studio

Editor visual offline para montar a statusline nativa do Codex. Escolha os itens disponíveis, ajuste a ordem, visualize uma prévia e gere TOML para copiar ou baixar.

## Uso

Abra [`index.html`](index.html) em um navegador. Não precisa de servidor, dependências ou acesso à internet.

O arquivo [`config.toml`](config.toml) traz a configuração padrão usada pelo editor. Para aplicá-la, incorpore a tabela `[tui]` ao seu `~/.codex/config.toml` ou ao caminho definido por `CODEX_HOME`, preservando as outras opções que já existirem. Reinicie o Codex para carregar a alteração.

A ordem padrão segue diretório, branch, modelo, contexto, tokens e cota.

## Custo

O Codex não aceita renderer customizado em `tui.status_line`, só identificadores nativos. Para custo, existem três caminhos:

- `estimated-thread-cost` e `thread-credits`: custo e créditos estimados da thread, **só em workspace Enterprise**. Fora disso o item fica oculto.
- `five-hour-limit` e `weekly-limit`: quanto resta da cota no Plus/Pro. Nesses planos é a métrica de gasto que importa, já que não há fatura por token.
- `used-tokens`, `total-input-tokens` e `total-output-tokens`: volume de tokens da sessão.

Histórico, gasto do dia/mês e gateway continuam exclusivos do `claude-statusline`.

## Itens disponíveis

O Studio lista todos os itens de `StatusLineItem` do TUI (`codex-rs/tui/src/bottom_pane/status_line_setup.rs`): diretório, projeto, branch, PR, mudanças da branch, modelo, raciocínio, fast mode, contexto (restante, usado, janela), tokens, limites de 5h e semanal, custo e créditos Enterprise, estado, progresso do plano, permissões, aprovação, thread, hostname, versão e afins. Itens sem dado ficam ocultos.

## Configuração padrão

```toml
[tui]
status_line = ["current-dir", "git-branch", "model-with-reasoning", "context-remaining", "used-tokens", "five-hour-limit", "weekly-limit", "estimated-thread-cost"]
terminal_title = ["spinner", "project"]
```

Confira a [referência oficial de configuração do Codex](https://developers.openai.com/codex/config-file/config-reference.md) para detalhes sobre `tui.status_line`.
