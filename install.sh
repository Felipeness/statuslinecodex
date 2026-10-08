#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "$0")" && pwd)"
BIN_DIR="$HOME/.local/bin"
WRAPPER="$BIN_DIR/codex-statusline"
ZSHRC="$HOME/.zshrc"

# Detecta binário real do Codex (node script, não o alias)
REAL_CODEX="$(mise which codex 2>/dev/null || \
  find "$HOME/.local/share/mise/installs/node" -name "codex.js" 2>/dev/null | head -1 || \
  command -v codex)"

echo "→ Compilando codex-statusline..."
cd "$REPO_DIR"
go build -o "$WRAPPER" .

echo "→ Binário em $WRAPPER"

# Atualiza alias no .zshrc
ALIAS_LINE='alias codex="CODEX_REAL_BIN='"$REAL_CODEX"' codex-statusline --dangerously-bypass-approvals-and-sandbox"'

if grep -q "codex-statusline" "$ZSHRC" 2>/dev/null; then
  # Substitui linha existente
  sed -i '' '/codex-statusline/d' "$ZSHRC"
fi

# Remove alias antigo simples se existir
sed -i '' '/^alias codex="codex --dangerously/d' "$ZSHRC" 2>/dev/null || true

echo "" >> "$ZSHRC"
echo "# codex-statusline wrapper" >> "$ZSHRC"
echo "$ALIAS_LINE" >> "$ZSHRC"

echo "→ Alias atualizado em $ZSHRC"
echo ""
echo "✓ Pronto. Recarregue o shell:"
echo "  source ~/.zshrc"
echo ""
echo "Teste: codex-statusline stats"
