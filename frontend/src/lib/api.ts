import { fetchAuthSession } from "aws-amplify/auth";
import type {
  CreateTokenRequest,
  CreateTokenResponse,
  GetProjectResponse,
  ListProjectsResponse,
  ListScansResponse,
  ListTokensResponse,
  TrendsResponse,
  VulnerabilitiesResponse,
} from "./types";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Custom error class for API errors
export class ApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
    this.name = "ApiError";
  }
}

// Get auth header from Cognito session
// Token is automatically refreshed by Amplify if expired
async function getAuthHeader(): Promise<string> {
  try {
    const session = await fetchAuthSession();
    const token = session.tokens?.idToken?.toString();
    if (!token) {
      throw new Error("No token available");
    }
    return `Bearer ${token}`;
  } catch {
    throw new ApiError(401, "Not authenticated");
  }
}

// Generic request function with auth
async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const authHeader = await getAuthHeader();

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: authHeader,
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new ApiError(response.status, error.error || "Request failed");
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}

// API Functions

export async function listProjects(): Promise<ListProjectsResponse> {
  return request("/api/projects");
}

export async function getProject(projectName: string): Promise<GetProjectResponse> {
  return request(`/api/projects/${encodeURIComponent(projectName)}`);
}

export async function listScans(projectName: string): Promise<ListScansResponse> {
  return request(`/api/projects/${encodeURIComponent(projectName)}/scans`);
}

export async function listTokens(): Promise<ListTokensResponse> {
  return request("/api/tokens");
}

export async function createToken(data: CreateTokenRequest): Promise<CreateTokenResponse> {
  return request("/api/tokens", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function deleteToken(tokenId: string): Promise<void> {
  await request(`/api/tokens/${encodeURIComponent(tokenId)}`, {
    method: "DELETE",
  });
}

export async function getTrends(projectId: string, days: number = 30): Promise<TrendsResponse> {
  return request(`/api/trends?projectId=${encodeURIComponent(projectId)}&days=${days}`);
}

export async function listVulnerabilities(
  scanId: string,
  options?: { severity?: string; packageName?: string; limit?: number }
): Promise<VulnerabilitiesResponse> {
  const params = new URLSearchParams({ scanId });
  if (options?.severity) params.set("severity", options.severity);
  if (options?.packageName) params.set("package", options.packageName);
  if (options?.limit) params.set("limit", options.limit.toString());

  return request(`/api/vulnerabilities?${params.toString()}`);
}
