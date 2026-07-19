import { useEffect, useState, type FormEvent } from "react";
import {
  fetchUnboundSettings,
  updateUnboundSettings,
  type UnboundSettings,
} from "../api/client";

export default function Unbound() {
  const [settings, setSettings] = useState<UnboundSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    const load = window.setTimeout(() => {
      fetchUnboundSettings()
        .then(setSettings)
        .catch((err: unknown) => {
          setError(err instanceof Error ? err.message : "Settings could not be loaded");
        })
        .finally(() => setLoading(false));
    }, 0);
    return () => clearTimeout(load);
  }, []);

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (!settings || saving) return;

    setSaving(true);
    setMessage("");
    setError("");
    try {
      const updated = await updateUnboundSettings(settings);
      setSettings(updated);
      setMessage("Configuration validated and Unbound restarted successfully.");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Configuration could not be saved");
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return <Page title="Unbound Settings"><p>Loading settings…</p></Page>;
  }
  if (!settings) {
    return <Page title="Unbound Settings"><p className="error-message">{error}</p></Page>;
  }

  return (
    <Page title="Unbound Settings">
      <p style={{ opacity: 0.7, maxWidth: "720px", marginBottom: "28px" }}>
        RootGuard validates every change inside the running resolver before it
        activates the configuration. Invalid values leave the active DNS
        configuration untouched.
      </p>

      <form onSubmit={submit} className="glass-card compact" style={{ maxWidth: "760px" }}>
        <Toggle
          label="QNAME minimisation"
          description="Reduces the amount of query information sent to upstream name servers."
          checked={settings.qname_minimisation}
          onChange={(value) => setSettings({ ...settings, qname_minimisation: value })}
        />
        <Toggle
          label="Prefetch popular records"
          description="Refreshes frequently used records before their cache entry expires."
          checked={settings.prefetch}
          onChange={(value) => setSettings({ ...settings, prefetch: value })}
        />
        <Toggle
          label="Serve expired records"
          description="Keeps DNS available during temporary upstream failures."
          checked={settings.serve_expired}
          onChange={(value) => setSettings({ ...settings, serve_expired: value })}
        />

        <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: "18px", marginTop: "26px" }}>
          <NumberField
            label="Minimum TTL"
            value={settings.cache_min_ttl}
            min={0}
            max={3600}
            onChange={(value) => setSettings({ ...settings, cache_min_ttl: value })}
          />
          <NumberField
            label="Maximum TTL"
            value={settings.cache_max_ttl}
            min={60}
            max={604800}
            onChange={(value) => setSettings({ ...settings, cache_max_ttl: value })}
          />
          <NumberField
            label="Resolver threads"
            value={settings.threads}
            min={1}
            max={32}
            onChange={(value) => setSettings({ ...settings, threads: value })}
          />
        </div>

        {message && <p style={{ color: "#22c55e", marginTop: "22px" }}>{message}</p>}
        {error && <p style={{ color: "#ef4444", marginTop: "22px" }}>{error}</p>}

        <button type="submit" disabled={saving} style={{ marginTop: "24px", padding: "12px 22px" }}>
          {saving ? "Validating…" : "Validate and apply"}
        </button>
      </form>
    </Page>
  );
}

function Page({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div style={{ padding: "40px 80px 60px", maxWidth: "1500px", margin: "0 auto" }}>
      <h1 style={{ fontSize: "28px", marginBottom: "18px" }}>{title}</h1>
      {children}
    </div>
  );
}

function Toggle({ label, description, checked, onChange }: {
  label: string;
  description: string;
  checked: boolean;
  onChange: (value: boolean) => void;
}) {
  return (
    <label style={{ display: "flex", justifyContent: "space-between", gap: "30px", padding: "16px 0", borderBottom: "1px solid rgba(255,255,255,.08)" }}>
      <span>
        <strong>{label}</strong>
        <span style={{ display: "block", opacity: 0.62, fontSize: "13px", marginTop: "5px" }}>{description}</span>
      </span>
      <input type="checkbox" checked={checked} onChange={(event) => onChange(event.target.checked)} />
    </label>
  );
}

function NumberField({ label, value, min, max, onChange }: {
  label: string;
  value: number;
  min: number;
  max: number;
  onChange: (value: number) => void;
}) {
  return (
    <label style={{ fontSize: "13px", opacity: 0.8 }}>
      {label}
      <input
        type="number"
        value={value}
        min={min}
        max={max}
        onChange={(event) => onChange(Number(event.target.value))}
        style={{ display: "block", width: "100%", marginTop: "8px", padding: "10px", boxSizing: "border-box" }}
      />
    </label>
  );
}
