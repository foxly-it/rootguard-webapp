import { useCallback, useEffect, useState, type FormEvent } from "react";
import {
  fetchUnboundDiagnostics,
  fetchUnboundHistory,
  fetchUnboundSettings,
  previewUnboundSettings,
  restoreUnboundVersion,
  updateUnboundSettings,
  type UnboundDiagnosticReport,
  type UnboundHistoryEntry,
  type UnboundPreview,
  type UnboundSettings,
} from "../api/client";
import "../styles/unbound.css";

export default function Unbound() {
  const [settings, setSettings] = useState<UnboundSettings | null>(null);
  const [history, setHistory] = useState<UnboundHistoryEntry[]>([]);
  const [preview, setPreview] = useState<UnboundPreview | null>(null);
  const [diagnostics, setDiagnostics] = useState<UnboundDiagnosticReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const reload = useCallback(async () => {
    const [loadedSettings, loadedHistory] = await Promise.all([
      fetchUnboundSettings(),
      fetchUnboundHistory(),
    ]);
    setSettings(loadedSettings);
    setHistory(loadedHistory);
  }, []);

  useEffect(() => {
    reload()
      .catch((err: unknown) => setError(errorMessage(err, "Unbound konnte nicht geladen werden")))
      .finally(() => setLoading(false));
  }, [reload]);

  async function createPreview(event: FormEvent) {
    event.preventDefault();
    if (!settings || busy) return;
    setBusy(true);
    clearFeedback();
    try {
      setPreview(await previewUnboundSettings(settings));
    } catch (err) {
      setError(errorMessage(err, "Vorschau konnte nicht erstellt werden"));
    } finally {
      setBusy(false);
    }
  }

  async function applyPreview() {
    if (!settings || busy || !preview?.changed) return;
    setBusy(true);
    clearFeedback();
    try {
      const updated = await updateUnboundSettings(settings);
      setSettings(updated);
      setPreview(null);
      await reload();
      setMessage("Konfiguration wurde validiert, versioniert und erfolgreich aktiviert.");
    } catch (err) {
      setError(errorMessage(err, "Konfiguration konnte nicht aktiviert werden"));
    } finally {
      setBusy(false);
    }
  }

  async function restore(entry: UnboundHistoryEntry) {
    if (busy || !window.confirm(`Version vom ${formatDate(entry.created_at)} wirklich wiederherstellen?`)) return;
    setBusy(true);
    clearFeedback();
    try {
      setSettings(await restoreUnboundVersion(entry.id));
      setPreview(null);
      await reload();
      setMessage("Die ausgewählte Version wurde validiert und wiederhergestellt.");
    } catch (err) {
      setError(errorMessage(err, "Version konnte nicht wiederhergestellt werden"));
    } finally {
      setBusy(false);
    }
  }

  async function runDiagnostics() {
    if (busy) return;
    setBusy(true);
    clearFeedback();
    try {
      setDiagnostics(await fetchUnboundDiagnostics());
    } catch (err) {
      setError(errorMessage(err, "Diagnose konnte nicht ausgeführt werden"));
    } finally {
      setBusy(false);
    }
  }

  function clearFeedback() {
    setMessage("");
    setError("");
  }

  if (loading) return <Page><p>Unbound-Einstellungen werden geladen…</p></Page>;
  if (!settings) return <Page><p className="error-message">{error}</p></Page>;

  return (
    <Page>
      <div className="unbound-heading">
        <div>
          <p className="unbound-eyebrow">RECURSIVE DNS</p>
          <h1>Unbound Configuration</h1>
          <p>Änderungen werden vor der Aktivierung im Resolver geprüft. RootGuard hält bis zu 20 funktionierende Versionen für einen sicheren Rollback bereit.</p>
        </div>
        <button className="secondary-action" type="button" disabled={busy} onClick={runDiagnostics}>
          {busy ? "Bitte warten…" : "Diagnose ausführen"}
        </button>
      </div>

      {message && <div className="feedback success">{message}</div>}
      {error && <div className="feedback error">{error}</div>}

      <div className="unbound-grid">
        <form onSubmit={createPreview} className="glass-card compact settings-panel">
          <h2>Resolver-Einstellungen</h2>
          <Toggle label="QNAME-Minimierung" description="Übermittelt autoritativen Nameservern nur den jeweils notwendigen Teil eines Namens." checked={settings.qname_minimisation} onChange={(value) => setSettings({ ...settings, qname_minimisation: value })} />
          <Toggle label="Beliebte Einträge vorladen" description="Aktualisiert häufig verwendete Cache-Einträge vor ihrem Ablauf." checked={settings.prefetch} onChange={(value) => setSettings({ ...settings, prefetch: value })} />
          <Toggle label="Abgelaufene Einträge bereitstellen" description="Hält DNS bei vorübergehenden Upstream-Problemen verfügbar." checked={settings.serve_expired} onChange={(value) => setSettings({ ...settings, serve_expired: value })} />

          <div className="number-grid">
            <NumberField label="Minimum TTL" value={settings.cache_min_ttl} min={0} max={3600} onChange={(value) => setSettings({ ...settings, cache_min_ttl: value })} />
            <NumberField label="Maximum TTL" value={settings.cache_max_ttl} min={60} max={604800} onChange={(value) => setSettings({ ...settings, cache_max_ttl: value })} />
            <NumberField label="Resolver-Threads" value={settings.threads} min={1} max={32} onChange={(value) => setSettings({ ...settings, threads: value })} />
          </div>

          <button type="submit" disabled={busy}>Änderungen prüfen</button>
        </form>

        <section className="glass-card compact side-panel">
          <h2>Diagnose</h2>
          {!diagnostics && <p className="muted-copy">Prüft Konfiguration, rekursive Auflösung und die Ablehnung ungültiger DNSSEC-Antworten.</p>}
          {diagnostics && (
            <>
              <div className={`overall-status ${diagnostics.healthy ? "healthy" : "failed"}`}>{diagnostics.healthy ? "Alle Prüfungen bestanden" : "Mindestens eine Prüfung fehlgeschlagen"}</div>
              {diagnostics.checks.map((check) => <DiagnosticRow key={check.name} {...check} />)}
              <small className="timestamp">Geprüft: {formatDate(diagnostics.checked_at)}</small>
            </>
          )}
        </section>
      </div>

      {preview && (
        <section className="glass-card preview-panel">
          <div className="panel-heading"><div><p className="unbound-eyebrow">VORSCHAU</p><h2>Konfigurationsänderungen</h2></div><button className="text-action" type="button" onClick={() => setPreview(null)}>Schließen</button></div>
          {!preview.changed ? <p>Die vorgeschlagene Konfiguration entspricht bereits dem aktiven Stand.</p> : (
            <>
              <div className="change-list">{preview.changes.map((change) => <div key={change.field}><code>{fieldLabel(change.field)}</code><span>{change.before}</span><b>→</b><span>{change.after}</span></div>)}</div>
              <details><summary>Generierte Unbound-Konfiguration anzeigen</summary><pre>{preview.rendered_config}</pre></details>
              <button type="button" disabled={busy} onClick={applyPreview}>{busy ? "Aktiviere…" : "Validieren und aktivieren"}</button>
            </>
          )}
        </section>
      )}

      <section className="glass-card history-panel">
        <div className="panel-heading"><div><p className="unbound-eyebrow">ROLLBACK</p><h2>Konfigurationsverlauf</h2></div><span>{history.length} / 20 Versionen</span></div>
        {history.length === 0 ? <p className="muted-copy">Nach der ersten Änderung erscheinen hier die Ausgangs- und die neue Version.</p> : (
          <div className="history-list">{history.map((entry, index) => (
            <article key={entry.id}>
              <div><strong>{index === 0 ? "Aktuellste Version" : "Gesicherte Version"}</strong><span>{formatDate(entry.created_at)}</span><small>Threads {entry.settings.threads} · TTL {entry.settings.cache_min_ttl}–{entry.settings.cache_max_ttl}</small></div>
              <button className="secondary-action" type="button" disabled={busy || index === 0} onClick={() => restore(entry)}>{index === 0 ? "Aktiv" : "Wiederherstellen"}</button>
            </article>
          ))}</div>
        )}
      </section>
    </Page>
  );
}

function Page({ children }: { children: React.ReactNode }) {
  return <div className="unbound-page">{children}</div>;
}

function Toggle({ label, description, checked, onChange }: { label: string; description: string; checked: boolean; onChange: (value: boolean) => void }) {
  return <label className="toggle-row"><span><strong>{label}</strong><small>{description}</small></span><input type="checkbox" checked={checked} onChange={(event) => onChange(event.target.checked)} /></label>;
}

function NumberField({ label, value, min, max, onChange }: { label: string; value: number; min: number; max: number; onChange: (value: number) => void }) {
  return <label className="number-field">{label}<input type="number" value={value} min={min} max={max} onChange={(event) => onChange(Number(event.target.value))} /></label>;
}

function DiagnosticRow({ name, passed, detail }: { name: string; passed: boolean; detail: string }) {
  return <div className="diagnostic-row"><span className={passed ? "check-pass" : "check-fail"}>{passed ? "✓" : "!"}</span><div><strong>{fieldLabel(name)}</strong><small>{detail}</small></div></div>;
}

function fieldLabel(field: string) {
  const labels: Record<string, string> = { qname_minimisation: "QNAME-Minimierung", prefetch: "Prefetch", serve_expired: "Serve Expired", cache_min_ttl: "Minimum TTL", cache_max_ttl: "Maximum TTL", threads: "Resolver-Threads", configuration: "Konfiguration", resolution: "DNS-Auflösung", dnssec: "DNSSEC" };
  return labels[field] ?? field;
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("de-DE", { dateStyle: "medium", timeStyle: "medium" }).format(new Date(value));
}

function errorMessage(error: unknown, fallback: string) {
  return error instanceof Error ? error.message : fallback;
}
