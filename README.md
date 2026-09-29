# Codex Statusline Studio

Editor visual offline para montar a statusline nativa do Codex. Escolha os itens disponíveis, ajuste a ordem, visualize uma prévia e gere TOML para copiar ou baixar.

## Uso

Abra [`index.html`](index.html) em um navegador. Não precisa de servidor, dependências ou acesso à internet.

O arquivo [`config.toml`](config.toml) traz a configuração padrão usada pelo editor. Para aplicá-la, incorpore a tabela `[tui]` ao seu `~/.codex/config.toml` ou ao caminho definido por `CODEX_HOME`, preservando as outras opções que já existirem. Reinicie o Codex para carregar a alteração.

A ordem padrão prioriza modelo e esforço, percentual de contexto restante e branch; o diretório vem por último para continuar visível em terminais estreitos.

## Itens disponíveis

- `model-with-reasoning`
- `context-remaining`
- `current-dir`
- `model`
- `git-branch`

`context-remaining` mostra o percentual de contexto restante. O TUI não expõe contagem exata de tokens ou custo como itens da statusline.

O Codex aceita identificadores nativos em `tui.status_line`; não oferece um renderer customizado de statusline. Por isso, o Studio configura e pré-visualiza os itens nativos, mas não implementa métricas de custo, gateway, histórico ou os componentes customizados do `claude-statusline`.

## Configuração padrão

```toml
[tui]
status_line = ["model-with-reasoning", "context-remaining", "git-branch", "current-dir"]
terminal_title = ["spinner", "project"]
```

Confira a [referência oficial de configuração do Codex](https://developers.openai.com/codex/config-file/config-reference.md) para detalhes sobre `tui.status_line`.
