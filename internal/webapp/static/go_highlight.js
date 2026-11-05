(() => {
  const keywords = new Set([
    "break", "case", "chan", "const", "continue", "default", "defer", "else", "fallthrough",
    "for", "func", "go", "goto", "if", "import", "interface", "map", "package", "range",
    "return", "select", "struct", "switch", "type", "var",
  ]);
  const types = new Set([
    "any", "bool", "byte", "comparable", "complex64", "complex128", "error", "float32", "float64",
    "int", "int8", "int16", "int32", "int64", "rune", "string", "uint", "uint8", "uint16",
    "uint32", "uint64", "uintptr",
  ]);
  const builtins = new Set([
    "append", "cap", "clear", "close", "complex", "copy", "delete", "imag", "len", "make",
    "max", "min", "new", "panic", "print", "println", "real", "recover",
  ]);
  const constants = new Set(["false", "iota", "nil", "true"]);
  const tokenPattern = /\/\/[^\n]*|\/\*[\s\S]*?\*\/|`[^`]*`|"(?:\\[\s\S]|[^"\\])*"|'(?:\\[\s\S]|[^'\\])*'|\b(?:0[xX][\da-fA-F_]+|0[bB][01_]+|0[oO][0-7_]+|(?:\d[\d_]*\.?[\d_]*|\.\d[\d_]*)(?:[eE][+-]?[\d_]+)?i?)\b|[A-Za-z_]\w*|[+\-*\/%&|^<>=!:~]+|[.;,(){}\[\]]/g;

  function escapeHTML(value) {
    return value.replace(/[&<>]/g, (character) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;" })[character]);
  }

  function classify(token, source, end) {
    if (token.startsWith("//") || token.startsWith("/*")) return "comment";
    if (token.startsWith('"') || token.startsWith("`") || token.startsWith("'")) return "string";
    if (/^(?:\d|\.\d)/.test(token)) return "number";
    if (keywords.has(token)) return "keyword";
    if (types.has(token)) return "type";
    if (builtins.has(token)) return "builtin";
    if (constants.has(token)) return "constant";
    if (/^[A-Za-z_]\w*$/.test(token) && /^\s*\(/.test(source.slice(end))) return "function";
    if (/^[+\-*\/%&|^<>=!:~]+$/.test(token)) return "operator";
    if (/^[.;,(){}\[\]]$/.test(token)) return "punctuation";
    return "";
  }

  function highlight(source) {
    let output = "";
    let cursor = 0;
    tokenPattern.lastIndex = 0;
    for (const match of source.matchAll(tokenPattern)) {
      output += escapeHTML(source.slice(cursor, match.index));
      const token = match[0];
      const className = classify(token, source, match.index + token.length);
      const escaped = escapeHTML(token);
      output += className ? `<span class="token ${className}">${escaped}</span>` : escaped;
      cursor = match.index + token.length;
    }
    return output + escapeHTML(source.slice(cursor));
  }

  window.GoHighlighter = Object.freeze({ highlight });
})();
