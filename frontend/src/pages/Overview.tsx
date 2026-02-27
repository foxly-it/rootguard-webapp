import { useEffect, useState } from "react";

export default function Overview() {

  const [apiStatus, setApiStatus] = useState("loading");

  useEffect(() => {
    fetch("/api/health")
      .then(res => res.json())
      .then(data => setApiStatus(data.status))
      .catch(() => setApiStatus("offline"));
  }, []);

  return (
    <>
      <h1>Dashboard Overview</h1>

      <div className="card-grid">
        <div className="card">
          <h3>Backend API</h3>
          <p className={apiStatus === "ok" ? "status-ok" : "status-error"}>
            {apiStatus}
          </p>
        </div>

        <div className="card">
          <h3>Docker</h3>
          <p>Checking...</p>
        </div>

        <div className="card">
          <h3>Unbound</h3>
          <p>Not configured</p>
        </div>

        <div className="card">
          <h3>AdGuard</h3>
          <p>Not configured</p>
        </div>
      </div>
    </>
  );
}