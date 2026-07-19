import { Routes, Route } from "react-router-dom";
import SidebarLayout from "./layout/SidebarLayout";
import Overview from "./pages/Overview";
import Unbound from "./pages/Unbound";
import AdGuard from "./pages/AdGuard";

function App() {
  return (
    <SidebarLayout>
      <Routes>
        <Route path="/" element={<Overview />} />
        <Route path="/unbound" element={<Unbound />} />
        <Route path="/adguard" element={<AdGuard />} />
      </Routes>
    </SidebarLayout>
  );
}

export default App;
