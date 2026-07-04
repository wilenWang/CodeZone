export type User = {
  id: number;
  username: string;
  displayName: string;
  userType: "human" | "agent";
};

export type LoginResult = {
  token: string;
  user: User;
};

export function authHeader(token: string | null): Record<string, string> {
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export async function apiGet<T>(path: string, token: string | null): Promise<T> {
  const response = await fetch(path, { headers: authHeader(token) });
  return parseResponse<T>(response);
}

export async function apiPost<T>(path: string, token: string | null, body: unknown): Promise<T> {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeader(token) },
    body: JSON.stringify(body),
  });
  return parseResponse<T>(response);
}

async function parseResponse<T>(response: Response): Promise<T> {
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message ?? "Request failed");
  }
  return data as T;
}

export function devLogin(username: string): Promise<LoginResult> {
  return apiPost<LoginResult>("/api/auth/dev-login", null, { username });
}
