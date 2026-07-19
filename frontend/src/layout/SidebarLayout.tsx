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
            to="/unbound"
            className={({ isActive }) =>
              isActive ? "nav-item active" : "nav-item"
            }
          >
            Unbound Settings
          </NavLink>

          <NavLink
            to="/adguard"
            className={({ isActive }) =>
              isActive ? "nav-item active" : "nav-item"
            }
          >
            AdGuard Home
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
