import { Routes, Route } from "react-router-dom";
import SidebarLayout from "./layout/SidebarLayout";
import Overview from "./pages/Overview";

function App() {
  return (
    <SidebarLayout>
      <Routes>
        <Route path="/" element={<Overview />} />
      </Routes>
    </SidebarLayout>
  );
}

export default App;