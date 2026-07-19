import { useEffect, useState } from "react";
import {
  bootstrapAdGuard,
  fetchAdGuardStatus,
  type AdGuardStatus,
} from "../api/client";

export default function AdGuard() {
  const [status, setStatus] = useState<AdGuardStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [bootstrapping, setBootstrapping] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    fetchAdGuardStatus()
      .then(setStatus)
      .catch((err: unknown) => setError(errorMessage(err, "AdGuard status could not be loaded")))
      .finally(() => setLoading(false));
  }, []);

  async function initialize() {
    if (bootstrapping) return;
    setBootstrapping(true);
    setMessage("");
    setError("");
    try {
      const updated = await bootstrapAdGuard();
      setStatus(updated);
      setMessage("AdGuard Home is configured and forwards DNS exclusively to Unbound.");
    } catch (err) {
      setError(errorMessage(err, "AdGuard could not be initialized"));
    } finally {
      setBootstrapping(false);
    }
  }

  return (
    <div style={{ padding: "40px 80px 60px", maxWidth: "1500px", margin: "0 auto" }}>
      <h1 style={{ fontSize: "28px", marginBottom: "18px" }}>AdGuard Home</h1>
      <p style={{ opacity: 0.7, maxWidth: "720px", marginBottom: "28px" }}>
        RootGuard performs the initial setup, stores the generated credentials inside
        the Core service and connects AdGuard Home to the validating Unbound resolver.
      </p>

      <section className="glass-card compact" style={{ maxWidth: "760px" }}>
        {loading ? <p>Checking AdGuard Home…</p> : (
          <>
            <StatusRow label="Installation" active={Boolean(status?.configured)} activeText="Configured" inactiveText="Setup required" />
            <StatusRow label="Service" active={Boolean(status?.healthy)} activeText="Healthy" inactiveText="Unavailable" />
            <StatusRow label="Unbound upstream" active={Boolean(status?.upstream_ready)} activeText="Validated" inactiveText="Not validated" />
            <div style={{ marginTop: "22px", opacity: 0.72, fontSize: "14px" }}>
              DNS upstream: <code>{status?.upstream || "rootguard-unbound:5335"}</code>
            </div>
          </>
        )}

        {message && <p style={{ color: "#22c55e", marginTop: "22px" }}>{message}</p>}
        {error && <p className="error-message" style={{ marginTop: "22px" }}>{error}</p>}

        {!loading && !status?.configured && (
          <button type="button" disabled={bootstrapping} onClick={initialize} style={{ marginTop: "24px", padding: "12px 22px" }}>
            {bootstrapping ? "Initializing…" : "Initialize securely"}
          </button>
        )}
      </section>
    </div>
  );
}

function StatusRow({ label, active, activeText, inactiveText }: {
  label: string;
  active: boolean;
  activeText: string;
  inactiveText: string;
}) {
  return (
    <div style={{ display: "flex", justifyContent: "space-between", padding: "16px 0", borderBottom: "1px solid rgba(255,255,255,.08)" }}>
      <strong>{label}</strong>
      <span style={{ color: active ? "#22c55e" : "#f59e0b" }}>{active ? activeText : inactiveText}</span>
    </div>
  );
}

function errorMessage(error: unknown, fallback: string) {
  return error instanceof Error ? error.message : fallback;
}
