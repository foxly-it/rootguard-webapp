import { Routes, Route } from "react-router-dom";
import SidebarLayout from "./layout/SidebarLayout";
import Overview from "./pages/Overview";
import Unbound from "./pages/Unbound";

function App() {
  return (
    <SidebarLayout>
      <Routes>
        <Route path="/" element={<Overview />} />
		<Route path="/unbound" element={<Unbound />} />
      </Routes>
    </SidebarLayout>
  );
}

export default App;
