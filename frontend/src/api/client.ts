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
    throw new Error(`API Error: ${response.status}`);
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
