"use client";

import { useQuery } from "@tanstack/react-query";
import { Header } from "@/components/layout/Header";
import { listProjects } from "@/lib/api";
import type { Project } from "@/lib/types";
import { FolderKanban } from "lucide-react";

export default function DashboardPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["projects"],
    queryFn: listProjects,
  });
  const projects = data?.projects ?? [];
  return (
    <>
      <Header title="Dashboard" />
      <div className="dashboard-content">
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-label">Total Projects</div>
            <div className="stat-value">{isLoading ? "-" : projects.length}</div>
          </div>
        </div>

        <div className="stat-card">
          <h2
            style={{ fontSize: "0.875rem", fontWeight: 500, color: "#fff", marginBottom: "1rem" }}
          >
            Recent Projects
          </h2>

          {isLoading ? (
            <div className="skeleton" style={{ height: "100px" }} />
          ) : error ? (
            <div style={{ color: "#fca5a5" }}>Failed to load</div>
          ) : projects.length === 0 ? (
            <div className="empty-state">
              <FolderKanban size={48} className="empty-state-icon" />
              <div className="empty-state-title">No projects yet</div>
              <div className="empty-state-description">Run a CLI scan to create a project</div>
              <code
                style={{
                  padding: "0.75rem",
                  background: "#111",
                  borderRadius: "6px",
                  color: "#a1a1aa",
                }}
              >
                vulnly scan ./
              </code>
            </div>
          ) : (
            <table className="data-table">
              <thead>
                <tr>
                  <th>Project</th>
                  <th>Last Scan</th>
                </tr>
              </thead>
              <tbody>
                {projects.slice(0, 5).map((p: Project) => (
                  <tr key={p.projectId}>
                    <td style={{ color: "#fff" }}>{p.name}</td>
                    <td>{p.lastScanAt ? new Date(p.lastScanAt).toLocaleDateString() : "Never"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </>
  );
}
