"use client";

import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { Header } from "@/components/layout/Header";
import { listProjects } from "@/lib/api";
import { getHealthStatus, getHealthLabel, getTotalVulns, timeAgo } from "@/lib/utils";
import type { ProjectWithStats } from "@/lib/types";
import { FolderKanban, Search, ArrowUpDown } from "lucide-react";
import "./projects.css";

export default function ProjectsPage() {
  const router = useRouter();
  const [search, setSearch] = useState("");
  const [sortBy, setSortBy] = useState<"name" | "lastScan" | "vulns">("lastScan");
  const [sortDir, setSortDir] = useState<"asc" | "desc">("desc");

  const { data, isLoading, error } = useQuery({
    queryKey: ["projects"],
    queryFn: listProjects,
  });

  const projects = useMemo(() => data?.projects ?? [], [data]);

  // Filter + sort
  const filtered = useMemo(() => {
    const result = projects.filter((p) => p.name.toLowerCase().includes(search.toLowerCase()));

    result.sort((a, b) => {
      let cmp = 0;
      if (sortBy === "name") {
        cmp = a.name.localeCompare(b.name);
      } else if (sortBy === "lastScan") {
        cmp = new Date(a.lastScanAt || 0).getTime() - new Date(b.lastScanAt || 0).getTime();
      } else if (sortBy === "vulns") {
        cmp = getTotalVulns(a.latestScan?.summary) - getTotalVulns(b.latestScan?.summary);
      }
      return sortDir === "asc" ? cmp : -cmp;
    });

    return result;
  }, [projects, search, sortBy, sortDir]);

  function handleSort(col: "name" | "lastScan" | "vulns") {
    if (sortBy === col) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortBy(col);
      setSortDir(col === "name" ? "asc" : "desc");
    }
  }

  function handleRowClick(project: ProjectWithStats) {
    router.push(`/projects/${encodeURIComponent(project.name)}`);
  }

  return (
    <>
      <Header title="Projects" />
      <div className="dashboard-content">
        {/* Stats row */}
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-label">Total Projects</div>
            <div className="stat-value">{isLoading ? "-" : projects.length}</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">With Vulnerabilities</div>
            <div className="stat-value">
              {isLoading
                ? "-"
                : projects.filter((p) => getTotalVulns(p.latestScan?.summary) > 0).length}
            </div>
          </div>
          <div className="stat-card">
            <div className="stat-label">Critical Issues</div>
            <div className="stat-value stat-value-critical">
              {isLoading
                ? "-"
                : projects.reduce((sum, p) => sum + (p.latestScan?.summary?.critical ?? 0), 0)}
            </div>
          </div>
        </div>

        {/* Search bar */}
        <div className="projects-search-bar">
          <Search size={16} className="search-icon" />
          <input
            type="text"
            className="search-input"
            placeholder="Search projects..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <span className="search-count">
            {filtered.length} project{filtered.length !== 1 ? "s" : ""}
          </span>
        </div>

        {/* Table */}
        {isLoading ? (
          <div className="stat-card">
            <div className="skeleton" style={{ height: "200px" }} />
          </div>
        ) : error ? (
          <div className="stat-card" style={{ color: "#fca5a5" }}>
            Failed to load projects
          </div>
        ) : filtered.length === 0 && !search ? (
          <div className="stat-card">
            <div className="empty-state">
              <FolderKanban size={48} className="empty-state-icon" />
              <div className="empty-state-title">No projects yet</div>
              <div className="empty-state-description">
                Run a CLI scan to create your first project
              </div>
              <code className="empty-code">vulnly scan ./</code>
            </div>
          </div>
        ) : filtered.length === 0 && search ? (
          <div className="stat-card">
            <div className="empty-state">
              <Search size={48} className="empty-state-icon" />
              <div className="empty-state-title">No matches</div>
              <div className="empty-state-description">
                No projects matching &quot;{search}&quot;
              </div>
            </div>
          </div>
        ) : (
          <div className="stat-card" style={{ padding: 0, overflow: "hidden" }}>
            <table className="data-table">
              <thead>
                <tr>
                  <th className="sortable-th" onClick={() => handleSort("name")}>
                    Project
                    {sortBy === "name" && <SortIndicator dir={sortDir} />}
                  </th>
                  <th>Ecosystem</th>
                  <th className="sortable-th" onClick={() => handleSort("lastScan")}>
                    Last Scan
                    {sortBy === "lastScan" && <SortIndicator dir={sortDir} />}
                  </th>
                  <th>Health</th>
                  <th className="sortable-th" onClick={() => handleSort("vulns")}>
                    Vulnerabilities
                    {sortBy === "vulns" && <SortIndicator dir={sortDir} />}
                  </th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((p) => {
                  const health = getHealthStatus(p.latestScan?.summary);
                  const total = getTotalVulns(p.latestScan?.summary);
                  return (
                    <tr
                      key={p.projectId}
                      className="clickable-row"
                      onClick={() => handleRowClick(p)}
                    >
                      <td>
                        <div className="project-name-cell">
                          <FolderKanban size={16} className="project-icon" />
                          <span className="project-name">{p.name}</span>
                        </div>
                      </td>
                      <td>{p.latestScan?.ecosystem ?? "-"}</td>
                      <td>{p.lastScanAt ? timeAgo(p.lastScanAt) : "Never"}</td>
                      <td>
                        <span className={`severity-badge ${health}`}>{getHealthLabel(health)}</span>
                      </td>
                      <td>
                        {total > 0 ? (
                          <VulnBar summary={p.latestScan!.summary} />
                        ) : (
                          <span style={{ color: "#71717a" }}>None</span>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </>
  );
}

// Sort direction indicator
function SortIndicator({ dir }: { dir: "asc" | "desc" }) {
  return (
    <ArrowUpDown
      size={12}
      style={{ marginLeft: 4, opacity: 0.7, transform: dir === "asc" ? "scaleY(-1)" : undefined }}
    />
  );
}

// Mini severity bar
function VulnBar({
  summary,
}: {
  summary: { critical: number; high: number; medium: number; low: number };
}) {
  const total = summary.critical + summary.high + summary.medium + summary.low;
  if (total === 0) return null;

  return (
    <div className="vuln-bar-container">
      <div className="vuln-bar">
        {summary.critical > 0 && (
          <div className="vuln-bar-segment critical" style={{ flex: summary.critical }} />
        )}
        {summary.high > 0 && (
          <div className="vuln-bar-segment high" style={{ flex: summary.high }} />
        )}
        {summary.medium > 0 && (
          <div className="vuln-bar-segment medium" style={{ flex: summary.medium }} />
        )}
        {summary.low > 0 && <div className="vuln-bar-segment low" style={{ flex: summary.low }} />}
      </div>
      <span className="vuln-bar-total">{total}</span>
    </div>
  );
}
