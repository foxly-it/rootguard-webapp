// =====================================================
// File: frontend/src/pages/Overview.tsx
// Purpose: Dashboard with Real Backend Hook
// =====================================================

import { useEffect, useState } from "react";

import CardHeader from "../components/CardHeader";
import MiniTrend from "../components/MiniTrend";
import ProgressBar from "../components/ProgressBar";
import ServiceControls from "../components/ServiceControls";

import "../styles/glass.css";

import type { ServiceStatus } from "../types/status";

import {
  fetchDashboard,
  serviceAction,
} from "../api/client";

// =====================================================
// Types
// =====================================================

interface DockerStats {
  cpu: number;
  memory: number;
  containers: number;
  status: ServiceStatus;
}

interface DnsStats {
  status: ServiceStatus;
  resolver: string;
  dnssec: boolean;
  trend: number[];
}

// =====================================================
// Component
// =====================================================

export default function Overview() {
  const [dockerStats, setDockerStats] = useState<DockerStats>({
    cpu: 0,
    memory: 0,
    containers: 0,
    status: "healthy",
  });

  const [dnsStats, setDnsStats] = useState<DnsStats>({
    status: "healthy",
    resolver: "",
    dnssec: false,
    trend: [12, 16, 14, 20, 24, 22, 30],
  });

  const [lastChecked, setLastChecked] = useState(
    new Date().toLocaleTimeString()
  );

  // =====================================================
  // Load Dashboard Data from Backend
  // =====================================================

  async function loadDashboard() {
    try {
      const data = await fetchDashboard();

      setDockerStats({
        cpu: data.docker.cpu,
        memory: data.docker.memory,
        containers: data.docker.containers,
        status: data.docker.status,
      });

      setDnsStats(prev => ({
        ...prev,
        status: data.dns.status,
        resolver: data.dns.resolver,
        dnssec: data.dns.dnssec,
      }));

      setLastChecked(new Date().toLocaleTimeString());
    } catch (err) {
      console.error("Dashboard load failed:", err);
    }
  }

  useEffect(() => {
    loadDashboard();

    const interval = setInterval(() => {
      loadDashboard();
    }, 10000);

    return () => clearInterval(interval);
  }, []);

  // =====================================================
  // Render
  // =====================================================

  return (
    <div
      style={{
        padding: "40px 80px 60px 80px",
        maxWidth: "1500px",
        margin: "0 auto",
      }}
    >
      <h1
        style={{
          fontSize: "28px",
          fontWeight: 600,
          marginBottom: "36px",
          opacity: 0.85,
        }}
      >
        System Overview
      </h1>

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(2, minmax(520px, 1fr))",
          gap: "48px",
        }}
      >
        {/* DNS CARD */}
        <div className="glass-card compact">
          <CardHeader
            title="DNS Engine"
            status={dnsStats.status}
            lastChecked={lastChecked}
          />

          <div style={{ marginBottom: "24px" }}>
            <MiniTrend data={dnsStats.trend} />
          </div>

          <InfoRow label="Resolver" value={dnsStats.resolver} />
          <InfoRow
            label="DNSSEC"
            value={dnsStats.dnssec ? "Enabled" : "Disabled"}
          />

          <div style={{ marginTop: "28px" }}>
            <ServiceControls
              onStart={() => serviceAction("dns", "start")}
              onStop={() => serviceAction("dns", "stop")}
              onRestart={() => serviceAction("dns", "restart")}
            />
          </div>
        </div>

        {/* DOCKER CARD */}
        <div className="glass-card compact">
          <CardHeader
            title="Docker Stack"
            status={dockerStats.status}
            lastChecked={lastChecked}
          />

          <InfoRow
            label="Containers"
            value={`${dockerStats.containers} Running`}
            highlight
          />

          <MetricBlock
            label="CPU Usage"
            value={`${dockerStats.cpu}%`}
          >
            <ProgressBar value={dockerStats.cpu} color="#22c55e" />
          </MetricBlock>

          <MetricBlock
            label="Memory Usage"
            value={`${dockerStats.memory}%`}
          >
            <ProgressBar value={dockerStats.memory} color="#3b82f6" />
          </MetricBlock>

          <div style={{ marginTop: "28px" }}>
            <ServiceControls
              onStart={() => serviceAction("docker", "start")}
              onStop={() => serviceAction("docker", "stop")}
              onRestart={() => serviceAction("docker", "restart")}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

// =====================================================
// Helper Components
// =====================================================

function InfoRow({
  label,
  value,
  highlight,
}: {
  label: string;
  value: string;
  highlight?: boolean;
}) {
  return (
    <div style={{ marginBottom: "18px" }}>
      <div style={{ fontSize: "13px", opacity: 0.6 }}>
        {label}
      </div>
      <div
        style={{
          fontSize: highlight ? "20px" : "16px",
          fontWeight: highlight ? 600 : 500,
        }}
      >
        {value}
      </div>
    </div>
  );
}

function MetricBlock({
  label,
  value,
  children,
}: {
  label: string;
  value: string;
  children: React.ReactNode;
}) {
  return (
    <div style={{ marginBottom: "22px" }}>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          fontSize: "14px",
          marginBottom: "8px",
        }}
      >
        <span style={{ opacity: 0.7 }}>{label}</span>
        <span style={{ fontWeight: 600 }}>{value}</span>
      </div>
      {children}
    </div>
  );
}