const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("gitaryao_token");
}

export function setToken(token: string) {
  localStorage.setItem("gitaryao_token", token);
}

export function removeToken() {
  localStorage.removeItem("gitaryao_token");
}

export function isAuthenticated(): boolean {
  return !!getToken();
}

async function fetchAPI<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...((options.headers as Record<string, string>) || {}),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(error.error || `API error: ${res.status}`);
  }

  return res.json();
}

// --- Auth ---

export interface User {
  id: number;
  username: string;
  email: string;
  avatar_url: string;
  created_at: string;
}

export function getGithubLoginURL(): string {
  return `${API_URL}/auth/github`;
}

export function getGoogleLoginURL(): string {
  return `${API_URL}/auth/google`;
}

export function getCloneURL(owner: string, repo: string): string {
  return `${API_URL}/${owner}/${repo}.git`;
}

export function getCurrentUser(): Promise<User> {
  return fetchAPI<User>("/api/user");
}

// --- Repos ---

export interface Repository {
  id: number;
  owner_id: number;
  owner_name: string;
  name: string;
  description: string;
  is_private: boolean;
  default_branch: string;
  created_at: string;
  updated_at: string;
}

export function listMyRepos(): Promise<Repository[]> {
  return fetchAPI<Repository[]>("/api/repos");
}

export function listPublicRepos(): Promise<Repository[]> {
  return fetchAPI<Repository[]>("/api/repos/public");
}

export function listUserRepos(username: string): Promise<Repository[]> {
  return fetchAPI<Repository[]>(`/api/users/${username}/repos`);
}

export function createRepo(data: {
  name: string;
  description: string;
  is_private: boolean;
}): Promise<Repository> {
  return fetchAPI<Repository>("/api/repos", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

// --- Browse ---

export interface TreeEntry {
  name: string;
  type: "blob" | "tree";
  size: number;
}

export interface CommitInfo {
  hash: string;
  message: string;
  author: string;
  email: string;
  date: string;
}

export interface BranchInfo {
  name: string;
  is_default: boolean;
}

export function getRepoTree(
  owner: string,
  repo: string,
  ref: string,
  path: string = ""
): Promise<TreeEntry[]> {
  const encodedPath = path ? `/${encodeURIComponent(path)}` : "";
  return fetchAPI<TreeEntry[]>(
    `/api/repos/${owner}/${repo}/tree/${ref}${encodedPath}`
  );
}

export function getRepoBlob(
  owner: string,
  repo: string,
  ref: string,
  path: string
): Promise<{ content: string; size: number }> {
  return fetchAPI(`/api/repos/${owner}/${repo}/blob/${ref}/${encodeURIComponent(path)}`);
}

export function getRepoCommits(
  owner: string,
  repo: string,
  ref: string
): Promise<CommitInfo[]> {
  return fetchAPI<CommitInfo[]>(`/api/repos/${owner}/${repo}/commits/${ref}`);
}

export function getRepoBranches(
  owner: string,
  repo: string
): Promise<BranchInfo[]> {
  return fetchAPI<BranchInfo[]>(`/api/repos/${owner}/${repo}/branches`);
}
