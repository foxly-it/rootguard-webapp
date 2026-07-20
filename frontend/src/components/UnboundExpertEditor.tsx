import { useCallback, useEffect, useMemo, useRef, useState, type KeyboardEvent, type UIEvent } from "react";
import {
  fetchUnboundCustom,
  fetchUnboundDirectives,
  previewUnboundCustom,
  updateUnboundCustom,
  type UnboundCustomPreview,
  type UnboundDirectiveReference,
} from "../api/client";
import "../styles/unbound-expert.css";

const templates = [
  { name: "Lokaler DNS-Eintrag", content: 'server:\n    local-zone: "home.arpa." static\n    local-data: "router.home.arpa. 300 IN A 192.168.1.1"\n' },
  { name: "Private Antwort erlauben", content: 'server:\n    private-domain: "home.arpa"\n' },
  { name: "Bedingte Weiterleitung", content: 'forward-zone:\n    name: "corp.example."\n    forward-addr: 192.168.1.53\n' },
  { name: "Zusätzliche Härtung", content: "server:\n    hide-identity: yes\n    hide-version: yes\n    harden-glue: yes\n    harden-dnssec-stripped: yes\n    aggressive-nsec: yes\n" },
];

export default function UnboundExpertEditor({ version, onActivated }: { version?: string; onActivated: () => Promise<void> }) {
  const [open, setOpen] = useState(false);
  const [active, setActive] = useState("");
  const [draft, setDraft] = useState("");
  const [maxBytes, setMaxBytes] = useState(65536);
  const [directives, setDirectives] = useState<UnboundDirectiveReference[]>([]);
  const [preview, setPreview] = useState<UnboundCustomPreview | null>(null);
  const [cursor, setCursor] = useState(0);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const highlightRef = useRef<HTMLPreElement>(null);

  const load = useCallback(async () => {
    const [document, reference] = await Promise.all([fetchUnboundCustom(), fetchUnboundDirectives()]);
    setActive(document.content);
    setDraft(document.content);
    setMaxBytes(document.max_bytes);
    setDirectives(reference);
    setPreview(null);
  }, []);

  useEffect(() => {
    load().catch((err: unknown) => setError(errorMessage(err, "Expertenkonfiguration konnte nicht geladen werden")));
  }, [load, version]);

  const prefix = useMemo(() => directivePrefix(draft, cursor), [draft, cursor]);
  const suggestions = useMemo(() => {
    if (!prefix) return [];
    return directives.filter((item) => item.name.toLowerCase().startsWith(prefix.toLowerCase())).slice(0, 6);
  }, [directives, prefix]);
  const selectedDirective = useMemo(() => directiveAt(draft, cursor, directives), [draft, cursor, directives]);
  const bytes = new TextEncoder().encode(draft).length;

  async function createPreview() {
    if (busy) return;
    setBusy(true);
    clearFeedback();
    try {
      const result = await previewUnboundCustom(draft);
      setDraft(result.content);
      setPreview(result);
      setMessage("Entwurf wurde durch RootGuard-Regeln und unbound-checkconf geprüft.");
    } catch (err) {
      setPreview(null);
      setError(errorMessage(err, "Expertenkonfiguration ist ungültig"));
    } finally {
      setBusy(false);
    }
  }

  async function activate() {
    if (!preview?.changed || busy || !window.confirm("Geprüfte Expertenkonfiguration aktivieren und Unbound neu starten?")) return;
    setBusy(true);
    clearFeedback();
    try {
      const document = await updateUnboundCustom(preview.content);
      setActive(document.content);
      setDraft(document.content);
      setPreview(null);
      await onActivated();
      setMessage("90-rootguard-custom.conf wurde versioniert und sicher aktiviert.");
    } catch (err) {
      setError(errorMessage(err, "Expertenkonfiguration konnte nicht aktiviert werden"));
    } finally {
      setBusy(false);
    }
  }

  function insertSuggestion(reference: UnboundDirectiveReference) {
    const start = Math.max(0, cursor - prefix.length);
    const indentation = draft.slice(draft.lastIndexOf("\n", start - 1) + 1, start).match(/^\s*/)?.[0] ?? "";
    const example = reference.example.includes("\n") ? reference.example : indentation + reference.example.trimStart();
    const next = draft.slice(0, start - indentation.length) + example + draft.slice(cursor);
    const nextCursor = start - indentation.length + example.length;
    setDraft(next);
    setPreview(null);
    requestAnimationFrame(() => {
      textareaRef.current?.focus();
      textareaRef.current?.setSelectionRange(nextCursor, nextCursor);
      setCursor(nextCursor);
    });
  }

  function insertTemplate(content: string) {
    const separator = draft && !draft.endsWith("\n\n") ? (draft.endsWith("\n") ? "\n" : "\n\n") : "";
    setDraft(draft + separator + content);
    setPreview(null);
  }

  function handleKeyDown(event: KeyboardEvent<HTMLTextAreaElement>) {
    if (event.key !== "Tab") return;
    event.preventDefault();
    const element = event.currentTarget;
    const next = draft.slice(0, element.selectionStart) + "    " + draft.slice(element.selectionEnd);
    const nextCursor = element.selectionStart + 4;
    setDraft(next);
    setPreview(null);
    requestAnimationFrame(() => element.setSelectionRange(nextCursor, nextCursor));
  }

  function syncScroll(event: UIEvent<HTMLTextAreaElement>) {
    if (!highlightRef.current) return;
    highlightRef.current.scrollTop = event.currentTarget.scrollTop;
    highlightRef.current.scrollLeft = event.currentTarget.scrollLeft;
  }

  function clearFeedback() {
    setError("");
    setMessage("");
  }

  return (
    <section className="glass-card expert-panel">
      <div className="panel-heading expert-heading">
        <div><p className="unbound-eyebrow">EXPERTENMODUS</p><h2>Zusätzliche Unbound-Konfiguration</h2></div>
        <button className="secondary-action" type="button" onClick={() => setOpen(!open)}>{open ? "Editor schließen" : "Editor öffnen"}</button>
      </div>
      <p className="muted-copy">Bearbeitet ausschließlich <code>90-rootguard-custom.conf</code>. Die unveränderliche Basis und die geführten Einstellungen bleiben getrennt.</p>
      {!open && <div className="expert-summary"><span>{active ? `${new TextEncoder().encode(active).length} Bytes aktiv` : "Keine Zusatzkonfiguration aktiv"}</span><span>Versioniert · Checkconf · Auto-Rollback</span></div>}
      {open && <>
        <div className="expert-warning"><strong>Expertenmodus</strong><span>Fehlerhafte DNS-Regeln können Namen unerreichbar machen. RootGuard blockiert systemkritische Direktiven und aktiviert nur geprüfte Entwürfe.</span></div>
        {message && <div className="feedback success">{message}</div>}
        {error && <div className="feedback error">{error}</div>}
        <div className="template-row">{templates.map((template) => <button type="button" key={template.name} onClick={() => insertTemplate(template.content)}>{template.name}</button>)}</div>
        <div className="expert-grid">
          <div>
            <div className="editor-toolbar"><span>90-rootguard-custom.conf</span><span className={bytes > maxBytes ? "limit-exceeded" : ""}>{bytes.toLocaleString("de-DE")} / {maxBytes.toLocaleString("de-DE")} Bytes</span></div>
            <div className="code-editor">
              <pre ref={highlightRef} aria-hidden="true">{highlightConfig(draft)}{"\n"}</pre>
              <textarea ref={textareaRef} aria-label="Zusätzliche Unbound-Konfiguration" spellCheck={false} value={draft} onChange={(event) => { setDraft(event.target.value); setCursor(event.target.selectionStart); setPreview(null); }} onSelect={(event) => setCursor(event.currentTarget.selectionStart)} onKeyDown={handleKeyDown} onScroll={syncScroll} />
            </div>
            {suggestions.length > 0 && <div className="completion-list"><span>Vervollständigen</span>{suggestions.map((item) => <button type="button" key={`${item.section}-${item.name}`} onClick={() => insertSuggestion(item)}><code>{item.name}</code><small>{item.section}</small></button>)}</div>}
            <div className="editor-actions">
              <button type="button" disabled={busy || bytes > maxBytes} onClick={createPreview}>{busy ? "Prüfe…" : "Entwurf prüfen"}</button>
              <button className="secondary-action" type="button" disabled={busy || draft === active} onClick={() => { setDraft(active); setPreview(null); }}>Aktiven Stand laden</button>
            </div>
          </div>
          <aside className="directive-help">
            <p className="unbound-eyebrow">KONTEXTHILFE</p>
            {selectedDirective ? <><div className="directive-title"><code>{selectedDirective.name}</code><span className={`risk-${selectedDirective.risk}`}>{riskLabel(selectedDirective.risk)}</span></div><p>{selectedDirective.description}</p><pre>{selectedDirective.example}</pre></> : <p className="muted-copy">Cursor auf eine bekannte Direktive setzen oder einige Buchstaben eingeben. RootGuard zeigt passende Vorschläge und deren Auswirkungen.</p>}
            <details><summary>Unterstützte Direktiven durchsuchen</summary><div className="directive-catalog">{directives.map((item) => <button type="button" key={`${item.section}-${item.name}`} onClick={() => insertSuggestion(item)}><code>{item.name}</code><small>{item.description}</small></button>)}</div></details>
          </aside>
        </div>
        {preview && <div className="custom-preview">
          <div className="validation-ok"><strong>✓ Syntaxprüfung bestanden</strong><span>{preview.validation}</span></div>
          <div className="custom-advice">{preview.advice.map((item) => <article className={`advice-item ${item.severity}`} key={item.id}><strong>{item.title}{item.line ? ` · Zeile ${item.line}` : ""}</strong><p>{item.description}</p><small>{item.suggestion}</small></article>)}</div>
          {!preview.changed ? <p>Der geprüfte Entwurf entspricht dem aktiven Stand.</p> : <button type="button" disabled={busy} onClick={activate}>Geprüften Entwurf aktivieren</button>}
        </div>}
      </>}
    </section>
  );
}

function directivePrefix(content: string, cursor: number) {
  return content.slice(0, cursor).match(/([a-z][a-z0-9-]*)$/i)?.[1] ?? "";
}

function directiveAt(content: string, cursor: number, directives: UnboundDirectiveReference[]) {
  const start = content.lastIndexOf("\n", Math.max(0, cursor - 1)) + 1;
  const end = content.indexOf("\n", cursor);
  const line = content.slice(start, end < 0 ? content.length : end).trim();
  const key = line.split(":", 1)[0].toLowerCase();
  return directives.find((item) => item.name.replace(/:$/, "").toLowerCase() === key);
}

function highlightConfig(content: string) {
  return content.split("\n").map((line, index) => {
    const commentAt = line.indexOf("#");
    const code = commentAt >= 0 ? line.slice(0, commentAt) : line;
    const comment = commentAt >= 0 ? line.slice(commentAt) : "";
    const separator = code.indexOf(":");
    if (separator < 0) return <span key={index}>{code}{comment && <span className="syntax-comment">{comment}</span>}{"\n"}</span>;
    const key = code.slice(0, separator);
    const value = code.slice(separator + 1);
    const valueClass = /\b(yes|no|allow|deny|refuse|static|transparent)\b/i.test(value) ? "syntax-keyword" : /\d/.test(value) ? "syntax-number" : "syntax-value";
    return <span key={index}>{key.slice(0, key.length - key.trimStart().length)}<span className="syntax-directive">{key.trim()}</span>:<span className={valueClass}>{value}</span>{comment && <span className="syntax-comment">{comment}</span>}{"\n"}</span>;
  });
}

function riskLabel(risk: string) {
  return ({ low: "Niedrig", medium: "Prüfen", high: "Hoch" } as Record<string, string>)[risk] ?? risk;
}

function errorMessage(error: unknown, fallback: string) {
  return error instanceof Error ? error.message : fallback;
}
