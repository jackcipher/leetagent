const state = {
  problems: [],
  current: null,
  reference: null,
  referenceLoading: false,
  referenceFile: null,
  progress: loadStoredJSON("leetagent.progress", {}),
  dirty: false,
  running: false,
};

const elements = {
  list: document.querySelector("#problemList"),
  workspace: document.querySelector("#workspace"),
  search: document.querySelector("#problemSearch"),
  sidebar: document.querySelector("#sidebar"),
  menuButton: document.querySelector("#menuButton"),
  description: document.querySelector("#description"),
  hints: document.querySelector("#hints"),
  reference: document.querySelector("#reference"),
  editor: document.querySelector("#sourceEditor"),
  sourceHighlight: document.querySelector("#sourceHighlight"),
  sourceHighlightCode: document.querySelector("#sourceHighlightCode"),
  sourceHighlightSentinel: document.querySelector("#sourceHighlightSentinel"),
  lineNumbers: document.querySelector("#lineNumbers"),
  workspaceResizer: document.querySelector("#workspaceResizer"),
  consoleResizer: document.querySelector("#consoleResizer"),
  cursorPosition: document.querySelector("#cursorPosition"),
  courseProgress: document.querySelector("#courseProgress"),
  problemPath: document.querySelector("#problemPath"),
  difficulty: document.querySelector("#difficulty"),
  saveButton: document.querySelector("#saveButton"),
  formatButton: document.querySelector("#formatButton"),
  runButton: document.querySelector("#runButton"),
  saveState: document.querySelector("#saveState"),
  testSummary: document.querySelector("#testSummary"),
  testOutput: document.querySelector("#testOutput"),
  toast: document.querySelector("#toast"),
};

async function request(url, options = {}) {
  const response = await fetch(url, {
    ...options,
    headers: { "Content-Type": "application/json", ...(options.headers || {}) },
  });
  const body = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(body.error || `Request failed with status ${response.status}`);
  }
  return body;
}

async function boot() {
  bindEvents();
  restoreLayout();
  try {
    const data = await request("/api/problems");
    state.problems = data.problems;
    renderProblemList();
    const requested = new URLSearchParams(location.search).get("problem");
    const initial = state.problems.find((problem) => problem.id === requested) || state.problems[0];
    if (initial) await selectProblem(initial.id);
    loadReference();
  } catch (error) {
    showToast(error.message);
    elements.description.innerHTML = `<h1>Unable to load course</h1><p>${escapeHTML(error.message)}</p>`;
  }
}

function bindEvents() {
  elements.search.addEventListener("input", renderProblemList);
  elements.menuButton.addEventListener("click", toggleSidebar);
  elements.editor.addEventListener("input", () => {
    state.dirty = true;
    elements.saveState.textContent = "Unsaved";
    updateLineNumbers();
    updateSourceHighlight();
    updateCursorPosition();
  });
  elements.editor.addEventListener("scroll", () => {
    elements.lineNumbers.scrollTop = elements.editor.scrollTop;
    elements.sourceHighlight.scrollTop = elements.editor.scrollTop;
    elements.sourceHighlight.scrollLeft = elements.editor.scrollLeft;
  });
  elements.editor.addEventListener("keydown", handleEditorKeydown);
  elements.editor.addEventListener("keyup", updateCursorPosition);
  elements.editor.addEventListener("click", updateCursorPosition);
  elements.editor.addEventListener("select", updateCursorPosition);
  elements.saveButton.addEventListener("click", saveSource);
  elements.formatButton.addEventListener("click", formatSource);
  elements.runButton.addEventListener("click", runTests);
  bindResizer(elements.workspaceResizer, resizeWorkspace);
  bindResizer(elements.consoleResizer, resizeConsole);
  document.querySelectorAll(".tab").forEach((tab) => {
    tab.addEventListener("click", () => activateTab(tab.dataset.tab));
  });
  window.addEventListener("beforeunload", (event) => {
    if (!state.dirty) return;
    event.preventDefault();
    event.returnValue = "";
  });
  document.addEventListener("keydown", (event) => {
    if ((event.metaKey || event.ctrlKey) && event.key === "s") {
      event.preventDefault();
      saveSource();
    }
    if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
      event.preventDefault();
      runTests();
    }
  });
}

function renderProblemList() {
  const query = elements.search.value.trim().toLowerCase();
  const visible = state.problems.filter((problem) => `${problem.title} ${problem.id}`.toLowerCase().includes(query));
  const passed = state.problems.filter((problem) => state.progress[problem.id] === "passed").length;
  elements.courseProgress.textContent = `${passed} / ${state.problems.length}`;
  const stages = new Map();
  visible.forEach((problem) => {
    if (!stages.has(problem.stage)) stages.set(problem.stage, []);
    stages.get(problem.stage).push(problem);
  });

  elements.list.replaceChildren();
  stages.forEach((problems) => {
    const heading = document.createElement("div");
    heading.className = "stage-title";
    const stagePassed = problems.filter((problem) => state.progress[problem.id] === "passed").length;
    heading.innerHTML = `<span>${escapeHTML(cleanStageTitle(problems[0].stageTitle))}</span><span class="stage-count">${stagePassed}/${problems.length}</span>`;
    elements.list.append(heading);
    problems.forEach((problem, index) => {
      const status = state.progress[problem.id] || "";
      const statusMark = status === "passed" ? "&#10003;" : status === "failed" ? "&#215;" : "";
      const button = document.createElement("button");
      button.type = "button";
      button.className = `problem-item${state.current?.id === problem.id ? " active" : ""}`;
      button.innerHTML = `<span class="problem-index">${String(index + 1).padStart(2, "0")}</span><span class="problem-name">${escapeHTML(cleanProblemTitle(problem.title))}</span><span class="problem-status ${status}">${statusMark}</span>`;
      button.title = problem.title;
      button.addEventListener("click", () => selectProblem(problem.id));
      elements.list.append(button);
    });
  });
}

async function selectProblem(id) {
  if (state.current?.id === id) {
    elements.sidebar.classList.remove("open");
    return;
  }
  if (state.dirty && !confirm("Discard unsaved changes and open another problem?")) return;

  setBusy(true);
  try {
    const problem = await request(`/api/problem?id=${encodeURIComponent(id)}`);
    state.current = problem;
    state.dirty = false;
    elements.problemPath.textContent = problem.title;
    elements.difficulty.textContent = problem.difficulty;
    elements.description.innerHTML = renderMarkdown(problem.statement);
    elements.hints.innerHTML = renderMarkdown(problem.hints || "# Hints\n\nNo hints are available for this exercise.");
    elements.editor.value = problem.source;
    elements.editor.selectionStart = 0;
    elements.editor.selectionEnd = 0;
    elements.editor.scrollTop = 0;
    elements.saveState.textContent = "Saved";
    elements.testSummary.textContent = "Not run";
    elements.testSummary.className = "test-summary";
    elements.testOutput.innerHTML = '<div class="test-empty">No results</div>';
    history.replaceState(null, "", `?problem=${encodeURIComponent(problem.id)}`);
    updateLineNumbers();
    updateSourceHighlight();
    updateCursorPosition();
    renderProblemList();
    elements.sidebar.classList.remove("open");
  } catch (error) {
    showToast(error.message);
  } finally {
    setBusy(false);
  }
}

async function saveSource() {
  if (!state.current || state.running) return;
  setBusy(true);
  elements.saveState.textContent = "Saving";
  try {
    await request(`/api/source?id=${encodeURIComponent(state.current.id)}`, {
      method: "PUT",
      body: JSON.stringify({ source: elements.editor.value }),
    });
    state.dirty = false;
    elements.saveState.textContent = "Saved";
    showToast("Source saved");
  } catch (error) {
    elements.saveState.textContent = "Save failed";
    showToast(error.message);
  } finally {
    setBusy(false);
  }
}

async function formatSource() {
  if (!state.current || state.running) return;
  setBusy(true);
  try {
    const result = await request("/api/format", {
      method: "POST",
      body: JSON.stringify({ source: elements.editor.value }),
    });
    if (result.source !== elements.editor.value) {
      const cursor = elements.editor.selectionStart;
      elements.editor.value = result.source;
      elements.editor.selectionStart = Math.min(cursor, result.source.length);
      elements.editor.selectionEnd = elements.editor.selectionStart;
      state.dirty = true;
      elements.saveState.textContent = "Unsaved";
      updateLineNumbers();
      updateSourceHighlight();
      updateCursorPosition();
    }
    showToast("Source formatted");
  } catch (error) {
    showToast(error.message);
  } finally {
    setBusy(false);
  }
}

async function runTests() {
  if (!state.current || state.running) return;
  state.running = true;
  setBusy(true);
  elements.testSummary.textContent = "Running";
  elements.testSummary.className = "test-summary";
  elements.testOutput.innerHTML = '<div class="test-running"><span class="test-spinner"></span><span>Running go test</span></div>';
  try {
    const result = await request(`/api/test?id=${encodeURIComponent(state.current.id)}`, {
      method: "POST",
      body: JSON.stringify({ source: elements.editor.value }),
    });
    state.dirty = false;
    elements.saveState.textContent = "Saved";
    setProblemProgress(state.current.id, result.passed ? "passed" : "failed");
    renderTestResult(result);
  } catch (error) {
    elements.testSummary.textContent = "Runner error";
    elements.testSummary.className = "test-summary failed";
    elements.testOutput.innerHTML = `<pre class="test-diagnostics">${escapeHTML(error.message)}</pre>`;
  } finally {
    state.running = false;
    setBusy(false);
  }
}

function activateTab(name) {
  document.querySelectorAll(".tab").forEach((tab) => {
    const active = tab.dataset.tab === name;
    tab.classList.toggle("active", active);
    tab.setAttribute("aria-selected", String(active));
  });
  document.querySelectorAll(".content-panel").forEach((pane) => pane.classList.toggle("active", pane.id === name));
  if (name === "reference" && !state.reference) loadReference();
}

async function loadReference() {
  if (state.reference || state.referenceLoading) return;
  state.referenceLoading = true;
  try {
    state.reference = await request("/api/reference");
    renderReference();
  } catch (error) {
    elements.reference.innerHTML = `<div class="reference-error"><h1>Common package unavailable</h1><p>${escapeHTML(error.message)}</p></div>`;
  } finally {
    state.referenceLoading = false;
  }
}

function renderReference() {
  const reference = state.reference;
  const groups = new Map();
  reference.symbols.forEach((symbol) => {
    const label = symbolGroupLabel(symbol.kind);
    if (!groups.has(label)) groups.set(label, []);
    groups.get(label).push(symbol);
  });
  const symbols = [...groups.entries()].map(([label, entries]) => `
    <section class="reference-group">
      <h2>${escapeHTML(label)}</h2>
      ${entries.map((symbol) => `
        <button class="reference-symbol" type="button" data-file="${escapeHTML(symbol.file)}" data-line="${symbol.line}">
          <span class="reference-symbol-meta"><span class="reference-kind">${escapeHTML(symbol.kind)}</span><span>${escapeHTML(symbol.file)}:${symbol.line}</span></span>
          <code>${escapeHTML(symbol.signature)}</code>
          ${symbol.doc ? `<span class="reference-doc">${escapeHTML(symbol.doc)}</span>` : ""}
        </button>`).join("")}
    </section>`).join("");

  elements.reference.innerHTML = `
    <header class="reference-header">
      <div>
        <h1>${escapeHTML(reference.importPath)}</h1>
        <p>${escapeHTML(reference.doc || "课程练习共用的类型、客户端和离线测试工具。")}</p>
      </div>
      <code>import &quot;${escapeHTML(reference.importPath)}&quot;</code>
    </header>
    <div class="reference-section-title">Exported API</div>
    <div class="reference-symbols">${symbols}</div>
    <div class="reference-source-header">
      <strong id="referenceFileName"></strong>
      <div class="reference-files" role="tablist">
        ${reference.files.map((file) => `<button class="reference-file" type="button" data-file="${escapeHTML(file.name)}">${escapeHTML(file.name)}</button>`).join("")}
      </div>
    </div>
    <div id="referenceSource" class="reference-source"></div>`;

  elements.reference.querySelectorAll(".reference-file").forEach((button) => {
    button.addEventListener("click", () => selectReferenceFile(button.dataset.file));
  });
  elements.reference.querySelectorAll(".reference-symbol").forEach((button) => {
    button.addEventListener("click", () => selectReferenceFile(button.dataset.file, Number(button.dataset.line)));
  });
  selectReferenceFile(state.referenceFile || reference.files[0].name);
}

function selectReferenceFile(filename, line = 0) {
  const file = state.reference?.files.find((entry) => entry.name === filename);
  if (!file) return;
  state.referenceFile = filename;
  elements.reference.querySelector("#referenceFileName").textContent = filename;
  elements.reference.querySelectorAll(".reference-file").forEach((button) => {
    const active = button.dataset.file === filename;
    button.classList.toggle("active", active);
    button.setAttribute("aria-selected", String(active));
  });
  const source = elements.reference.querySelector("#referenceSource");
  source.innerHTML = file.source.split("\n").map((content, index) => `
    <div id="reference-line-${index + 1}" class="reference-code-line">
      <span>${index + 1}</span><code>${escapeHTML(content) || " "}</code>
    </div>`).join("");
  if (line > 0) {
    requestAnimationFrame(() => {
      const target = elements.reference.querySelector(`#reference-line-${line}`);
      if (!target) return;
      target.classList.add("highlight");
      target.scrollIntoView({ behavior: "smooth", block: "center" });
      setTimeout(() => target.classList.remove("highlight"), 1800);
    });
  }
}

function symbolGroupLabel(kind) {
  if (kind === "interface") return "Interfaces";
  if (kind === "type") return "Types";
  if (kind === "var" || kind === "const") return "Variables and constants";
  if (kind === "func") return "Functions";
  if (kind === "method") return "Methods";
  return "Other";
}

function setBusy(busy) {
  elements.saveButton.disabled = busy;
  elements.formatButton.disabled = busy;
  elements.runButton.disabled = busy;
}

function updateLineNumbers() {
  const count = elements.editor.value.split("\n").length;
  elements.lineNumbers.textContent = Array.from({ length: count }, (_, index) => index + 1).join("\n");
}

function updateSourceHighlight() {
  const source = elements.editor.value;
  const highlighted = window.GoHighlighter?.highlight(source) || escapeHTML(source);
  elements.sourceHighlightCode.innerHTML = highlighted;
  elements.sourceHighlightSentinel.textContent = source.endsWith("\n") ? " " : "";
  elements.sourceHighlight.scrollTop = elements.editor.scrollTop;
  elements.sourceHighlight.scrollLeft = elements.editor.scrollLeft;
}

function updateCursorPosition() {
  const before = elements.editor.value.slice(0, elements.editor.selectionStart);
  const lines = before.split("\n");
  elements.cursorPosition.textContent = `Ln ${lines.length}, Col ${lines.at(-1).length + 1}`;
}

function handleEditorKeydown(event) {
  if (event.key !== "Tab") return;
  event.preventDefault();
  const start = elements.editor.selectionStart;
  const end = elements.editor.selectionEnd;
  elements.editor.setRangeText("\t", start, end, "end");
  elements.editor.dispatchEvent(new Event("input"));
}

function formatDuration(milliseconds) {
  if (milliseconds < 1000) return `${milliseconds} ms`;
  return `${(milliseconds / 1000).toFixed(1)} s`;
}

function renderTestResult(result) {
  const report = parseTestReport(result.output);
  const passed = report.tests.filter((test) => test.action === "pass").length;
  const failed = report.tests.filter((test) => test.action === "fail").length;
  const duration = formatDuration(result.durationMs);
  elements.testSummary.textContent = result.timedOut ? `Timed out · ${duration}` : result.passed ? `${passed || 1} passed · ${duration}` : `${failed || 1} failed · ${duration}`;
  elements.testSummary.className = `test-summary ${result.passed ? "passed" : "failed"}`;

  const cases = report.tests.map((test) => {
    const action = test.action || "fail";
    const icon = action === "pass" ? "&#10003;" : action === "skip" ? "-" : "&#215;";
    const elapsed = test.elapsed ? formatDuration(test.elapsed * 1000) : "";
    const details = cleanTestCaseOutput(test.output);
    if (action === "fail" || (action === "skip" && details)) {
      return `<details class="test-case ${action}" ${action === "fail" ? "open" : ""}>
        <summary><span class="test-case-icon">${icon}</span><span class="test-case-name">${escapeHTML(test.name)}</span><span class="test-case-time">${escapeHTML(elapsed)}</span></summary>
        ${details ? `<pre>${escapeHTML(details)}</pre>` : ""}
      </details>`;
    }
    return `<div class="test-case ${action}"><div class="test-case-row"><span class="test-case-icon">${icon}</span><span class="test-case-name">${escapeHTML(test.name)}</span><span class="test-case-time">${escapeHTML(elapsed)}</span></div></div>`;
  }).join("");

  const diagnostics = report.general.trim();
  const fallback = result.passed ? '<div class="test-case pass"><div class="test-case-row"><span class="test-case-icon">&#10003;</span><span class="test-case-name">Package passed</span><span class="test-case-time"></span></div></div>' : "";
  const diagnosticBlock = diagnostics && (!report.tests.length || !result.passed) ? `<pre class="test-diagnostics">${escapeHTML(diagnostics)}</pre>` : "";
  const raw = readableTestOutput(result.output);
  const rawBlock = raw ? `<details class="raw-output"><summary>Raw output</summary><pre>${escapeHTML(raw)}</pre></details>` : "";
  elements.testOutput.innerHTML = `${cases || fallback}${diagnosticBlock}${rawBlock}`;
}

function parseTestReport(output) {
  const tests = new Map();
  const general = [];
  for (const line of output.trim().split("\n")) {
    if (!line) continue;
    let event;
    try {
      event = JSON.parse(line);
    } catch {
      general.push(line);
      continue;
    }
    if (!event.Test) {
      if (event.Output) general.push(event.Output.replace(/\n$/, ""));
      continue;
    }
    if (!tests.has(event.Test)) tests.set(event.Test, { name: event.Test, action: "", elapsed: 0, output: [] });
    const test = tests.get(event.Test);
    if (event.Output) test.output.push(event.Output.replace(/\n$/, ""));
    if (["pass", "fail", "skip"].includes(event.Action)) {
      test.action = event.Action;
      test.elapsed = event.Elapsed || 0;
    }
  }
  const allTests = [...tests.values()];
  const visibleTests = allTests.filter((test) => {
    const hasSubtests = allTests.some((candidate) => candidate.name.startsWith(`${test.name}/`));
    return !hasSubtests || cleanTestCaseOutput(test.output) !== "";
  });
  return { tests: visibleTests, general: general.join("\n").replace(/\n{3,}/g, "\n\n") };
}

function cleanTestCaseOutput(lines) {
  return lines.filter((line) => !/^=== RUN\s/.test(line) && !/^--- (PASS|FAIL|SKIP):\s/.test(line)).join("\n").trim();
}

function readableTestOutput(output) {
  const readable = [];
  for (const line of output.trim().split("\n")) {
    if (!line) continue;
    try {
      const event = JSON.parse(line);
      if (event.Output) readable.push(event.Output.replace(/\n$/, ""));
    } catch {
      readable.push(line);
    }
  }
  return readable.join("\n").replace(/\n{3,}/g, "\n\n");
}

function setProblemProgress(id, status) {
  state.progress[id] = status;
  storeJSON("leetagent.progress", state.progress);
  renderProblemList();
}

function toggleSidebar() {
  if (window.matchMedia("(max-width: 780px)").matches) {
    elements.sidebar.classList.toggle("open");
    return;
  }
  elements.workspace.classList.toggle("sidebar-collapsed");
  const layout = loadStoredJSON("leetagent.layout", {});
  layout.sidebarCollapsed = elements.workspace.classList.contains("sidebar-collapsed");
  storeJSON("leetagent.layout", layout);
}

function bindResizer(element, onMove) {
  element.addEventListener("pointerdown", (event) => {
    if (window.matchMedia("(max-width: 780px)").matches) return;
    event.preventDefault();
    element.classList.add("dragging");
    const move = (moveEvent) => onMove(moveEvent);
    const stop = () => {
      element.classList.remove("dragging");
      document.removeEventListener("pointermove", move);
      document.removeEventListener("pointerup", stop);
    };
    document.addEventListener("pointermove", move);
    document.addEventListener("pointerup", stop);
  });
}

function resizeWorkspace(event) {
  const workspace = elements.workspace.getBoundingClientRect();
  const sidebarWidth = elements.workspace.classList.contains("sidebar-collapsed") ? 0 : elements.sidebar.getBoundingClientRect().width;
  const maximum = Math.max(320, workspace.width - sidebarWidth - 430);
  const width = Math.min(maximum, Math.max(320, event.clientX - workspace.left - sidebarWidth));
  document.documentElement.style.setProperty("--content-width", `${Math.round(width)}px`);
  const layout = loadStoredJSON("leetagent.layout", {});
  layout.contentWidth = Math.round(width);
  storeJSON("leetagent.layout", layout);
}

function resizeConsole(event) {
  const pane = document.querySelector(".work-pane").getBoundingClientRect();
  const maximum = Math.max(140, pane.height - 300);
  const height = Math.min(maximum, Math.max(140, pane.bottom - event.clientY - 38));
  document.documentElement.style.setProperty("--console-height", `${Math.round(height)}px`);
  const layout = loadStoredJSON("leetagent.layout", {});
  layout.consoleHeight = Math.round(height);
  storeJSON("leetagent.layout", layout);
}

function restoreLayout() {
  const layout = loadStoredJSON("leetagent.layout", {});
  if (Number.isFinite(layout.contentWidth)) document.documentElement.style.setProperty("--content-width", `${layout.contentWidth}px`);
  if (Number.isFinite(layout.consoleHeight)) document.documentElement.style.setProperty("--console-height", `${layout.consoleHeight}px`);
  if (layout.sidebarCollapsed) elements.workspace.classList.add("sidebar-collapsed");
}

function loadStoredJSON(key, fallback) {
  try {
    const value = JSON.parse(localStorage.getItem(key));
    return value && typeof value === "object" ? value : fallback;
  } catch {
    return fallback;
  }
}

function storeJSON(key, value) {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch {
    // Local persistence is optional; the editor remains usable without it.
  }
}

function renderMarkdown(markdown) {
  const blocks = [];
  let text = markdown.replace(/\r\n/g, "\n").replace(/```([^\n]*)\n([\s\S]*?)```/g, (_, language, code) => {
    const token = `@@CODE${blocks.length}@@`;
    blocks.push(`<pre><code data-language="${escapeHTML(language.trim())}">${escapeHTML(code.replace(/\n$/, ""))}</code></pre>`);
    return token;
  });
  text = escapeHTML(text);
  text = renderTables(text);

  const lines = text.split("\n");
  const output = [];
  let listType = null;
  let paragraph = [];

  const flushParagraph = () => {
    if (!paragraph.length) return;
    output.push(`<p>${inlineMarkdown(paragraph.join(" "))}</p>`);
    paragraph = [];
  };
  const closeList = () => {
    if (!listType) return;
    output.push(`</${listType}>`);
    listType = null;
  };

  for (const line of lines) {
    if (/^@@CODE\d+@@$/.test(line) || /^@@TABLE\d+@@$/.test(line)) {
      flushParagraph(); closeList(); output.push(line); continue;
    }
    const heading = line.match(/^(#{1,3})\s+(.+)$/);
    if (heading) {
      flushParagraph(); closeList();
      const level = heading[1].length;
      output.push(`<h${level}>${inlineMarkdown(heading[2])}</h${level}>`);
      continue;
    }
    const unordered = line.match(/^\s*[-*]\s+(.+)$/);
    const ordered = line.match(/^\s*\d+[.)、]\s+(.+)$/);
    if (unordered || ordered) {
      flushParagraph();
      const nextType = ordered ? "ol" : "ul";
      if (listType !== nextType) { closeList(); output.push(`<${nextType}>`); listType = nextType; }
      output.push(`<li>${inlineMarkdown((unordered || ordered)[1])}</li>`);
      continue;
    }
    if (line.startsWith("&gt; ")) {
      flushParagraph(); closeList(); output.push(`<blockquote><p>${inlineMarkdown(line.slice(5))}</p></blockquote>`); continue;
    }
    if (!line.trim()) { flushParagraph(); closeList(); continue; }
    paragraph.push(line.trim());
  }
  flushParagraph(); closeList();

  let html = output.join("\n");
  blocks.forEach((block, index) => {
    html = html.replace(`@@CODE${index}@@`, block).replace(`@@TABLE${index}@@`, block);
  });
  return html;

  function renderTables(value) {
    const rows = value.split("\n");
    const rendered = [];
    for (let index = 0; index < rows.length; index++) {
      const divider = /^\s*\|?\s*:?-{3,}:?\s*(\|\s*:?-{3,}:?\s*)+\|?\s*$/;
      if (index + 1 >= rows.length || !rows[index].includes("|") || !divider.test(rows[index + 1])) {
        rendered.push(rows[index]);
        continue;
      }
      const tableRows = [rows[index]];
      index += 2;
      while (index < rows.length && rows[index].includes("|")) {
        tableRows.push(rows[index]);
        index += 1;
      }
      index -= 1;
      const parsed = tableRows.map((row) => row.replace(/^\s*\||\|\s*$/g, "").split("|").map((cell) => inlineMarkdown(cell.trim())));
      const header = `<thead><tr>${parsed[0].map((cell) => `<th>${cell}</th>`).join("")}</tr></thead>`;
      const body = `<tbody>${parsed.slice(1).map((row) => `<tr>${row.map((cell) => `<td>${cell}</td>`).join("")}</tr>`).join("")}</tbody>`;
      const token = `@@TABLE${blocks.length}@@`;
      blocks.push(`<table>${header}${body}</table>`);
      rendered.push(token);
    }
    return rendered.join("\n");
  }
}

function inlineMarkdown(value) {
  return value
    .replace(/`([^`]+)`/g, "<code>$1</code>")
    .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
    .replace(/\[([^\]]+)]\(([^)]+)\)/g, '<span title="$2">$1</span>');
}

function cleanStageTitle(title) {
  return title.replace(/^阶段\s*\d+[：:]?\s*/, "").replace(/^Capstone[：:]?\s*/i, "Capstone ");
}

function cleanProblemTitle(title) {
  return title.replace(/^(ex\d+|M\d+)\s*/i, "").replace(/[（(](Easy|Medium|Hard)[）)]/i, "").trim();
}

function escapeHTML(value) {
  return String(value).replace(/[&<>"']/g, (character) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;" })[character]);
}

let toastTimer;
function showToast(message) {
  clearTimeout(toastTimer);
  elements.toast.textContent = message;
  elements.toast.classList.add("visible");
  toastTimer = setTimeout(() => elements.toast.classList.remove("visible"), 2400);
}

boot();
