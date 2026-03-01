// =====================================================
// File: frontend/src/components/StatusIndicator.tsx
// Purpose: Professional animated status indicator
// =====================================================

interface Props {
  status: "ok" | "warning" | "error";
}

export default function StatusIndicator({ status }: Props) {

  const label =
    status === "ok"
      ? "Operational"
      : status === "warning"
      ? "Degraded"
      : "Offline";

  return (
    <div className={`rg-status ${status}`}>
      <span className="rg-status-dot" />
      <span className="rg-status-label">{label}</span>
    </div>
  );
}
