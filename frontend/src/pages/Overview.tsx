// =====================================================
// File: frontend/src/pages/Overview.tsx
// Purpose: Clean Enterprise Dashboard Layout
// =====================================================

import Card from "../components/Card";
import StatusIndicator from "../components/StatusIndicator";

import DnsIcon from "../components/icons/DnsIcon";
import DockerIcon from "../components/icons/DockerIcon";

export default function Overview() {

  return (
    <div className="content-container">

      <h1 className="page-title">
        System Overview
      </h1>

      <div className="dashboard-grid">

        {/* ================= DNS ENGINE ================= */}
        <Card>

          <div className="card-header">
            <div className="card-title">
              <DnsIcon />
              <h3>DNS Engine</h3>
            </div>
            <StatusIndicator status="ok" />
          </div>

          <div className="card-body">

            <div className="card-info">
              <span>Resolver: Unbound</span>
              <span>DNSSEC: Enabled</span>
            </div>

            <div className="card-divider" />

            <div className="metrics-grid">

              <div className="metric">
                <div className="metric-value">24.5k</div>
                <div className="metric-label">Queries Today</div>
              </div>

              <div className="metric">
                <div className="metric-value">92%</div>
                <div className="metric-label">Cache Hit Rate</div>
              </div>

              <div className="metric">
                <div className="metric-value">12ms</div>
                <div className="metric-label">Avg Response</div>
              </div>

            </div>

          </div>
        </Card>

        {/* ================= DOCKER STACK ================= */}
        <Card>

          <div className="card-header">
            <div className="card-title">
              <DockerIcon />
              <h3>Docker Stack</h3>
            </div>
            <StatusIndicator status="ok" />
          </div>

          <div className="card-body">
            <div className="card-info">
              <span>AdGuard Home: Running</span>
              <span>Unbound: Running</span>
            </div>
          </div>

        </Card>

        {/* ================= CONFIGURATION ================= */}
        <Card className="card-full">

          <div className="card-header">
            <h3>Configuration</h3>
          </div>

          <div className="card-body">
            <div className="card-info">
              <span>Last Change: 2026-03-01</span>
              <span>Version: 0.1.0</span>
            </div>
          </div>

        </Card>

      </div>
    </div>
  );
}