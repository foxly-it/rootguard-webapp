// =====================================================
// File: frontend/src/layout/SidebarLayout.tsx
// Purpose: Layout with Sidebar Navigation
// =====================================================

import type { ReactNode } from "react";
import { NavLink } from "react-router-dom";
import Header from "./Header";

interface Props {
  children: ReactNode;
}

export default function SidebarLayout({ children }: Props) {

  return (
    <>
      <Header />

      <div className="layout">

        {/* ================= Sidebar ================= */}
        <aside className="sidebar">

          <NavLink
            to="/"
            className={({ isActive }) =>
              isActive ? "nav-item active" : "nav-item"
            }
          >
            Overview
          </NavLink>

          <NavLink
            to="/docker"
            className={({ isActive }) =>
              isActive ? "nav-item active" : "nav-item"
            }
          >
            Docker Stack
          </NavLink>

          <NavLink
            to="/assistant"
            className={({ isActive }) =>
              isActive ? "nav-item active" : "nav-item"
            }
          >
            Unbound Assistant
          </NavLink>

          <NavLink
            to="/health"
            className={({ isActive }) =>
              isActive ? "nav-item active" : "nav-item"
            }
          >
            Health
          </NavLink>

        </aside>

        {/* ================= Main ================= */}
        <main className="main">
          {children}
        </main>

      </div>
    </>
  );
}