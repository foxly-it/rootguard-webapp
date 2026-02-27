import { ReactNode } from "react";
import { Link } from "react-router-dom";

interface Props {
  children: ReactNode;
}

export default function SidebarLayout({ children }: Props) {
  return (
    <div className="layout">
      <aside className="sidebar">
        <h2>🦊 RootGuard</h2>

        <nav>
          <Link to="/">Overview</Link>
          <Link to="/docker">Docker Stack</Link>
          <Link to="/assistant">Unbound Assistant</Link>
          <Link to="/health">Health</Link>
        </nav>
      </aside>

      <main className="main">{children}</main>
    </div>
  );
}