// =====================================================
// File: frontend/src/api/client.ts
// Purpose: Central API client for RootGuard
// =====================================================

/**
 * Generic API helper
 */
async function request<T>(
  url: string,
  options?: RequestInit
): Promise<T> {
  const response = await fetch(url, {
    headers: {
      "Content-Type": "application/json",
    },
    ...options,
  });

  if (!response.ok) {
    const detail = (await response.text()).trim();
    throw new Error(detail || `API Error: ${response.status}`);
  }

  return response.json();
}

// =====================================================
// Dashboard Endpoint
// =====================================================

export interface DashboardResponse {
  docker: {
    cpu: number;
    memory: number;
    containers: number;
    status: "healthy" | "degraded" | "down";
  };
  dns: {
    status: "healthy" | "degraded" | "down";
    resolver: string;
    dnssec: boolean;
  };
}

export async function fetchDashboard() {
  return request<DashboardResponse>("/api/dashboard");
}

// =====================================================
// Service Action Endpoint
// =====================================================

export interface ServiceActionResponse {
  service: string;
  action: string;
  status: string;
}

export async function serviceAction(
  name: string,
  action: "start" | "stop" | "restart"
): Promise<void> {
  await request<ServiceActionResponse>(
    `/api/service/${name}/${action}`,
    {
      method: "POST",
    }
  );
}

// =====================================================
// Unbound Settings
// =====================================================

export interface UnboundSettings {
  qname_minimisation: boolean;
  prefetch: boolean;
  serve_expired: boolean;
  cache_min_ttl: number;
  cache_max_ttl: number;
  threads: number;
}

export async function fetchUnboundSettings(): Promise<UnboundSettings> {
  return request<UnboundSettings>("/api/unbound/settings");
}

export async function updateUnboundSettings(
  settings: UnboundSettings
): Promise<UnboundSettings> {
  return request<UnboundSettings>("/api/unbound/settings", {
    method: "PUT",
    body: JSON.stringify(settings),
  });
}

export interface UnboundChange {
  field: string;
  before: string;
  after: string;
}

export interface UnboundPreview {
  changed: boolean;
  changes: UnboundChange[];
  rendered_config: string;
}

export interface UnboundHistoryEntry {
  id: string;
  created_at: string;
  settings: UnboundSettings;
  config?: string;
  custom_config?: string;
}

export interface UnboundDiagnosticCheck {
  name: string;
  passed: boolean;
  detail: string;
}

export interface UnboundDiagnosticReport {
  healthy: boolean;
  checked_at: string;
  checks: UnboundDiagnosticCheck[];
}

export interface UnboundPreset {
  id: string;
  name: string;
  description: string;
  best_for: string;
  settings: UnboundSettings;
}

export interface UnboundRecommendation {
  id: string;
  severity: "success" | "recommendation" | "warning";
  field?: string;
  title: string;
  description: string;
  suggestion: string;
}

export interface UnboundAdvice {
  status: "optimized" | "suggestions" | "review";
  recommendations: UnboundRecommendation[];
}

export async function previewUnboundSettings(settings: UnboundSettings): Promise<UnboundPreview> {
  return request<UnboundPreview>("/api/unbound/preview", {
    method: "POST",
    body: JSON.stringify(settings),
  });
}

export async function fetchUnboundHistory(): Promise<UnboundHistoryEntry[]> {
  return request<UnboundHistoryEntry[]>("/api/unbound/history");
}

export async function restoreUnboundVersion(id: string): Promise<UnboundSettings> {
  return request<UnboundSettings>(`/api/unbound/history/${encodeURIComponent(id)}/restore`, {
    method: "POST",
  });
}

export async function fetchUnboundDiagnostics(): Promise<UnboundDiagnosticReport> {
  return request<UnboundDiagnosticReport>("/api/unbound/diagnostics");
}

export async function fetchUnboundPresets(): Promise<UnboundPreset[]> {
  return request<UnboundPreset[]>("/api/unbound/presets");
}

export async function fetchUnboundAdvice(settings: UnboundSettings): Promise<UnboundAdvice> {
  return request<UnboundAdvice>("/api/unbound/advice", {
    method: "POST",
    body: JSON.stringify(settings),
  });
}

export interface UnboundCustomDocument {
  content: string;
  max_bytes: number;
}

export interface UnboundCustomAdvice {
  id: string;
  severity: "success" | "recommendation" | "warning";
  line?: number;
  title: string;
  description: string;
  suggestion: string;
}

export interface UnboundCustomPreview {
  changed: boolean;
  content: string;
  validation: string;
  advice: UnboundCustomAdvice[];
}

export interface UnboundDirectiveReference {
  name: string;
  section: string;
  example: string;
  description: string;
  risk: "low" | "medium" | "high";
}

export async function fetchUnboundCustom(): Promise<UnboundCustomDocument> {
  return request<UnboundCustomDocument>("/api/unbound/custom");
}

export async function previewUnboundCustom(content: string): Promise<UnboundCustomPreview> {
  return request<UnboundCustomPreview>("/api/unbound/custom/preview", {
    method: "POST",
    body: JSON.stringify({ content }),
  });
}

export async function updateUnboundCustom(content: string): Promise<UnboundCustomDocument> {
  return request<UnboundCustomDocument>("/api/unbound/custom", {
    method: "PUT",
    body: JSON.stringify({ content }),
  });
}

export async function fetchUnboundDirectives(): Promise<UnboundDirectiveReference[]> {
  return request<UnboundDirectiveReference[]>("/api/unbound/directives");
}

export interface AdGuardStatus {
  configured: boolean;
  healthy: boolean;
  upstream: string;
  upstream_ready: boolean;
}

export async function fetchAdGuardStatus(): Promise<AdGuardStatus> {
  return request<AdGuardStatus>("/api/adguard/status");
}

export async function bootstrapAdGuard(): Promise<AdGuardStatus> {
  return request<AdGuardStatus>("/api/adguard/bootstrap", { method: "POST" });
}
