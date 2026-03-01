// =====================================================
// File: frontend/src/layout/Header.tsx
// Project: RootGuard WebApp
// Purpose: Global application header
//
// Behaviour:
// - Version wird immer angezeigt (auch "dev")
// - Commit wird nur angezeigt, wenn != "unknown"
// - Kein "vdev · unknown"
// - GitHub + Docs als Icon + Text
// - Icons folgen currentColor (Theme-ready)
// =====================================================

import { useEffect, useState } from "react";
import { getVersion } from "../services/api";

// RootGuard shield icon (SVG, icon-only)
import logo from "../assets/rootguard-icon.svg";

// Inline SVG icon components
import GithubIcon from "../components/icons/GithubIcon";
import DocsIcon from "../components/icons/DocsIcon";

interface VersionData {
  version: string;
  commit: string;
}

export default function Header() {
  const [version, setVersion] = useState<VersionData | null>(null);

  // -----------------------------------------------------
  // Fetch backend version metadata on mount
  // -----------------------------------------------------
  useEffect(() => {
    getVersion()
      .then(setVersion)
      .catch(() => setVersion(null));
  }, []);

  // -----------------------------------------------------
  // Format version string cleanly
  // - Always show version
  // - Only append commit if it is valid
  // -----------------------------------------------------
  function formatVersion(v: VersionData | null): string | null {
    if (!v) return null;

    if (v.commit && v.commit !== "unknown") {
      return `v${v.version} · ${v.commit.substring(0, 7)}`;
    }

    return `v${v.version}`;
  }

  return (
    <header className="rg-header">
      {/* =====================================================
          Left: Branding
         ===================================================== */}
      <div className="rg-header-left">
        <img
          src={logo}
          alt="RootGuard"
          className="rg-logo"
          aria-label="RootGuard Logo"
        />

        <h2 className="rg-title">RootGuard</h2>
      </div>

      {/* =====================================================
          Right: Version + Utility Links
         ===================================================== */}
      <div className="rg-header-right">

        {version && (
          <span className="rg-version">
            {formatVersion(version)}
          </span>
        )}

        <a
          href="https://github.com/foxly-it/rootguard-webapp"
          target="_blank"
          rel="noreferrer"
          className="rg-link"
          aria-label="GitHub Repository"
        >
          <GithubIcon />
          <span>GitHub</span>
        </a>

        <a
          href="https://foxly.de"
          target="_blank"
          rel="noreferrer"
          className="rg-link"
          aria-label="Documentation"
        >
          <DocsIcon />
          <span>Docs</span>
        </a>

      </div>
    </header>
  );
}