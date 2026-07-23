import { useCallback, useEffect, useMemo, useState } from "react";
import {
  Archive,
  CheckCircle2,
  Cpu,
  Download,
  LoaderCircle,
  RefreshCw,
  RotateCcw,
  ServerCog,
  ShieldCheck,
  PanelsTopLeft,
} from "lucide-react";
import {
  checkControlPlaneUpdates,
  checkUpdates,
  fetchControlPlaneUpdateStatus,
  fetchServices,
  fetchUpdateStatus,
  installServiceUpdate,
  installControlPlaneUpdates,
  serviceAction,
  type ServiceInfo,
  type UpdateServiceStatus,
  type UpdateStatus,
  type ControlPlaneUpdateStatus,
} from "../api/client";
import "../styles/stack.css";
import { useI18n } from "../i18n";

export default function Stack() {
  const { t, formatDate } = useI18n();
  const [updates, setUpdates] = useState<UpdateStatus | null>(null);
  const [controlPlane, setControlPlane] = useState<ControlPlaneUpdateStatus | null>(null);
  const [services, setServices] = useState<ServiceInfo[]>([]);
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    try {
      const [nextUpdates, nextControlPlane, nextServices] = await Promise.all([
        fetchUpdateStatus(),
        fetchControlPlaneUpdateStatus(),
        fetchServices(),
      ]);
      setUpdates(nextUpdates);
      setControlPlane(nextControlPlane);
      setServices(nextServices);
      setError("");
    } catch (cause) {
      setError(errorMessage(cause, "Stack-Status konnte nicht geladen werden."));
    }
  }, []);

  useEffect(() => {
    const initial = window.setTimeout(load, 0);
    return () => window.clearTimeout(initial);
  }, [load]);

  const busy = updates?.state === "checking" || updates?.state === "updating"
    || controlPlane?.state === "checking" || controlPlane?.state === "updating";
  useEffect(() => {
    if (!busy) return;
    const timer = window.setInterval(load, 1500);
    return () => window.clearInterval(timer);
  }, [busy, load]);

  const available = useMemo(
    () => (updates?.services.filter((service) => service.update_available).length ?? 0)
      + (controlPlane?.services.filter((service) => service.update_available).length ?? 0),
    [updates, controlPlane],
  );

  async function startCheck() {
    setError("");
    try {
      const [nextUpdates, nextControlPlane] = await Promise.all([
        checkUpdates(),
        checkControlPlaneUpdates(),
      ]);
      setUpdates(nextUpdates);
      setControlPlane(nextControlPlane);
    } catch (cause) {
      setError(errorMessage(cause, "Update-Prüfung konnte nicht gestartet werden."));
    }
  }

  async function startControlPlaneUpdate() {
    if (!window.confirm(t("stack.controlPlaneConfirm"))) return;
    setError("");
    try {
      setControlPlane(await installControlPlaneUpdates());
    } catch (cause) {
      setError(errorMessage(cause, t("stack.controlPlaneStartError")));
    }
  }

  async function startUpdate(service: UpdateServiceStatus) {
    const accepted = window.confirm(
      t("stack.confirmUpdate", { service: service.display_name }),
    );
    if (!accepted) return;
    setError("");
    try {
      setUpdates(await installServiceUpdate(service.name));
    } catch (cause) {
      setError(errorMessage(cause, "Update konnte nicht gestartet werden."));
    }
  }

  async function control(name: ServiceInfo["name"], action: "start" | "stop" | "restart") {
    if (action === "stop" && !window.confirm(t("stack.confirmStop", { service: name }))) return;
    try {
      await serviceAction(name, action);
      await load();
    } catch (cause) {
      setError(errorMessage(cause, "Dienstaktion fehlgeschlagen."));
    }
  }

  return (
    <div className="stack-page">
      <section className="stack-hero">
        <div>
          <span className="stack-eyebrow">CONTROLLED LIFECYCLE</span>
          <h1>{t("stack.title")}</h1>
          <p>{t("stack.intro")}</p>
        </div>
        <button type="button" className="stack-check-button" disabled={busy} onClick={startCheck}>
          {updates?.state === "checking" ? <LoaderCircle className="spin" size={17} /> : <RefreshCw size={17} />}
          {updates?.state === "checking" ? t("stack.checking") : t("stack.check")}
        </button>
      </section>

      {error && <div className="stack-feedback error">{error}</div>}
      {updates && (
        <div className={`stack-feedback ${updates.state === "failed" ? "error" : busy ? "working" : "success"}`}>
          {busy && <LoaderCircle className="spin" size={17} />}
          {!busy && updates.state !== "failed" && <CheckCircle2 size={17} />}
          <span>{updates.message}</span>
        </div>
      )}

      <section className="stack-summary">
        <Summary icon={<ServerCog />} label={t("stack.managed")} value={t("stack.runningCount", { count: services.filter((s) => s.status === "running").length })} />
        <Summary icon={<Download />} label={t("stack.available")} value={available ? t("stack.found", { count: available }) : t("stack.none")} />
        <Summary icon={<Archive />} label={t("stack.protection")} value={t("stack.backup")} />
        <Summary icon={<RotateCcw />} label={t("stack.failure")} value={t("stack.rollback")} />
      </section>

      {controlPlane && (
        <section className="control-plane-panel">
          <div className="control-plane-heading">
            <div>
              <span className="stack-eyebrow">{t("stack.controlPlaneEyebrow")}</span>
              <h2>{t("stack.controlPlaneTitle")}</h2>
              <p>{t("stack.controlPlaneIntro")}</p>
            </div>
            <button
              type="button"
              className="stack-check-button"
              disabled={busy || !controlPlane.services.some((service) => service.update_available)}
              onClick={startControlPlaneUpdate}
            >
              {controlPlane.state === "updating" ? <LoaderCircle className="spin" size={17} /> : <Download size={17} />}
              {controlPlane.state === "updating" ? t("stack.controlPlaneUpdating") : t("stack.controlPlaneInstall")}
            </button>
          </div>
          <div className={`control-plane-status ${controlPlane.state}`}>
            {controlPlane.state === "checking" || controlPlane.state === "updating"
              ? <LoaderCircle className="spin" size={15} />
              : <ShieldCheck size={15} />}
            <span>{controlPlane.message}</span>
          </div>
          <div className="control-plane-services">
            {controlPlane.services.map((service) => (
              <article key={service.name}>
                <span>{service.name === "core" ? <Cpu /> : <PanelsTopLeft />}</span>
                <div>
                  <strong>{service.display_name}</strong>
                  <small>{service.current_image || t("stack.notInspected")}</small>
                </div>
                <em className={service.update_available ? "available" : ""}>
                  {service.update_available ? t("stack.update") : service.checked_at ? t("stack.current") : t("stack.unchecked")}
                </em>
              </article>
            ))}
          </div>
          <p className="control-plane-note"><RotateCcw size={15} /> {t("stack.controlPlaneRollback")}</p>
        </section>
      )}

      <section className="stack-services">
        {updates?.services.map((service) => {
          const runtime = services.find((item) => item.name === service.name);
          const active = updates.active_service === service.name;
          return (
            <article className="stack-service-card" key={service.name}>
              <div className="stack-card-heading">
                <div>
                  <span className={`stack-runtime ${runtime?.status === "running" ? "running" : ""}`}>
                    {runtime?.status === "running" ? t("stack.running") : t("stack.stopped")}
                  </span>
                  <h2>{service.display_name}</h2>
                </div>
                <span className={`stack-update-badge ${service.update_available ? "available" : ""}`}>
                  {active ? t("stack.working") : service.update_available ? t("stack.update") : service.checked_at ? t("stack.current") : t("stack.unchecked")}
                </span>
              </div>

              <dl className="stack-image-data">
                <div><dt>{t("stack.runningImage")}</dt><dd>{service.current_image || "–"}</dd></div>
                <div><dt>{t("stack.target")}</dt><dd>{service.target_image}</dd></div>
                <div><dt>{t("stack.imageId")}</dt><dd>{shortID(service.current_id)}</dd></div>
                <div><dt>{t("stack.lastCheck")}</dt><dd>{service.checked_at && !service.checked_at.startsWith("0001-") ? formatDate(service.checked_at) : t("stack.neverChecked")}</dd></div>
              </dl>

              {service.error && <p className="stack-service-error">{service.error}</p>}

              <div className="stack-card-actions">
                <button type="button" disabled={busy || !service.update_available} onClick={() => startUpdate(service)}>
                  <Download size={15} /> {t("stack.install")}
                </button>
                <button type="button" disabled={busy} onClick={() => control(service.name, "restart")}>{t("common.restart")}</button>
                <button type="button" className="quiet" disabled={busy} onClick={() => control(service.name, runtime?.status === "running" ? "stop" : "start")}>
                  {runtime?.status === "running" ? t("common.stop") : t("common.start")}
                </button>
              </div>
            </article>
          );
        })}
      </section>

      <section className="stack-safety">
        <ShieldCheck size={23} />
        <div>
          <strong>{t("stack.safetyTitle")}</strong>
          <p>{t("stack.safetyText")}</p>
        </div>
      </section>
    </div>
  );
}

function Summary({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return <article><span>{icon}</span><div><small>{label}</small><strong>{value}</strong></div></article>;
}

function shortID(value?: string) {
  if (!value) return "–";
  return value.replace("sha256:", "").slice(0, 12);
}

function errorMessage(error: unknown, fallback: string) {
  return error instanceof Error ? error.message : fallback;
}
