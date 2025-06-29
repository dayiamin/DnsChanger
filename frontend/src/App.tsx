import { useEffect, useState } from "react";
import "./App.css";
import {
  GetDNSList,
  SetDNS,
  GetActiveDNS,
  PingDNS
} from "../wailsjs/go/main/App";

function App() {
  const [dnsList, setDnsList] = useState<Record<string, string[]>>({});
  const [active, setActive] = useState("");
  const [ping, setPing] = useState("-");
  const [pinging, setPinging] = useState(false);

  useEffect(() => {
    refresh();
  }, []);

  const refresh = async () => {
    const list = await GetDNSList();
    const act = await GetActiveDNS();
    setDnsList(list);
    setActive(act);
  };

  const handleSetDNS = async (name: string) => {
    const result = await SetDNS(name);
    alert(result);
    refresh();
  };

  const handleStartPing = async () => {
    setPinging(true);
    const result = await PingDNS();
    setPing(result);
    setPinging(false);
  };

  return (
    <div style={styles.app}>
      <h1 style={styles.title}>DNS Manager</h1>

      <div style={styles.grid}>
        {Object.entries(dnsList).map(([name, ips]) => (
          <div
            key={name}
            style={{
              ...styles.card,
              backgroundColor: name === active ? "#e0fff4" : "#ffffff",
              borderColor: name === active ? "#10b981" : "#ccc",
            }}
          >
            <div style={styles.cardContent}>
              <strong style={styles.dnsName}>{name}</strong>
              <div style={styles.ip}>{ips[0]}</div>
              <div style={styles.ip}>{ips[1]}</div>
              <button style={styles.button} onClick={() => handleSetDNS(name)}>
                Set This DNS
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* 👇 Ping Section */}
      <div style={styles.pingBox}>
        <button onClick={handleStartPing} style={styles.pingButton} disabled={pinging}>
          {pinging ? "Pinging..." : "Start Ping"}
        </button>
        <div style={{ marginTop: "10px" }}>{ping}</div>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  app: {
    fontFamily: "sans-serif",
    padding: "20px",
    backgroundColor: "#f1f5f9",
    minHeight: "100vh",
  },
  title: {
    textAlign: "center",
    marginBottom: "20px",
    color: "#1f2937",
  },
  grid: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))",
    gap: "16px",
    maxWidth: "960px",
    margin: "0 auto",
  },
  card: {
    border: "2px solid #ccc",
    borderRadius: "10px",
    padding: "20px",
    display: "flex",
    justifyContent: "center",
    transition: "all 0.2s ease-in-out",
    boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
    backgroundColor: "#fff",
  },
  cardContent: {
    textAlign: "center",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    gap: "8px",
  },
  dnsName: {
    fontSize: "18px",
    color: "#111827",
  },
  ip: {
    fontSize: "14px",
    color: "#4b5563",
  },
  button: {
    marginTop: "10px",
    padding: "8px 16px",
    backgroundColor: "#10b981",
    color: "#fff",
    border: "none",
    borderRadius: "6px",
    cursor: "pointer",
  },
  pingBox: {
    textAlign: "center",
    marginTop: "40px",
    fontSize: "14px",
    color: "#6b7280",
  },
  pingButton: {
    padding: "8px 16px",
    borderRadius: "6px",
    border: "none",
    backgroundColor: "#3b82f6",
    color: "white",
    cursor: "pointer",
  },
};

export default App;
