import type { VulnSummary } from "./types";

export type HealthStatus = "critical" | "high" | "medium" | "low" | "healthy";

export function getHealthStatus(summary?: VulnSummary): HealthStatus {
  if (!summary) return "healthy";
  if (summary.critical > 0) return "critical";
  if (summary.high > 0) return "high";
  if (summary.medium > 0) return "medium";
  if (summary.low > 0) return "low";
  return "healthy";
}

export function getHealthLabel(status: HealthStatus): string {
  if (status === "healthy") return "Healthy";
  return status.charAt(0).toUpperCase() + status.slice(1);
}

export function getTotalVulns(summary?: VulnSummary): number {
  if (!summary) return 0;
  return summary.critical + summary.high + summary.medium + summary.low;
}

export function timeAgo(dateString: string): string {
  const now = new Date();
  const date = new Date(dateString);
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (seconds < 60) return "just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  if (seconds < 2592000) return `${Math.floor(seconds / 86400)}d ago`;
  return date.toLocaleDateString();
}
