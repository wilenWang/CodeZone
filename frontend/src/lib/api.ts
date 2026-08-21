export type User = {
  id: number;
  username: string;
  displayName: string;
  avatarUrl: string | null;
  userType: "human" | "agent";
};

export type LoginResult = {
  token: string;
  user: User;
};

export type Conversation = {
  id: number;
  type: "direct" | "group";
  title: string | null;
  lastMessageId: number | null;
  lastMessageAt: string | null;
  unreadCount: number;
};

export type Message = {
  id: number;
  conversationId: number;
  senderId: number;
  contentMarkdown: string;
  contentPlain: string;
  createdAt: string;
};

export type AdminRow = Record<string, string | number | null>;

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

export async function apiPatch<T>(path: string, token: string | null, body: unknown): Promise<T> {
  const response = await fetch(path, {
    method: "PATCH",
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

export function updateProfile(token: string, displayName: string): Promise<User> {
  return apiPatch<User>("/api/me", token, { displayName });
}

export function uploadAvatar(token: string, file: File): Promise<User> {
  const form = new FormData();
  form.append("file", file);
  return fetch("/api/me/avatar", { method: "POST", headers: authHeader(token), body: form }).then(parseResponse<User>);
}

export function listConversations(token: string): Promise<{ conversations: Conversation[] }> {
  return apiGet("/api/conversations", token);
}

export function listUsers(token: string): Promise<{ users: User[] }> {
  return apiGet("/api/users", token);
}

export function ensureDirectConversation(token: string, userId: number): Promise<Conversation> {
  return apiPost<Conversation>("/api/conversations/direct", token, { userId });
}

export function createConversation(
  token: string,
  input: { type: "group"; title: string; memberIds: number[] },
): Promise<Conversation> {
  return apiPost<Conversation>("/api/conversations", token, input);
}

export function listMessages(token: string, conversationId: number): Promise<{ messages: Message[] }> {
  return apiGet(`/api/conversations/${conversationId}/messages?limit=50`, token);
}

export function sendMessage(token: string, conversationId: number, contentMarkdown: string): Promise<Message> {
  return apiPost(`/api/conversations/${conversationId}/messages`, token, { contentMarkdown });
}

export function markRead(token: string, conversationId: number): Promise<{ ok: boolean }> {
  return apiPost(`/api/conversations/${conversationId}/read`, token, {});
}

export function adminUsers(token: string): Promise<{ users: AdminRow[] }> {
  return apiGet("/api/admin/users", token);
}

export function adminConversations(token: string): Promise<{ conversations: AdminRow[] }> {
  return apiGet("/api/admin/conversations", token);
}

export function adminMessages(token: string): Promise<{ messages: AdminRow[] }> {
  return apiGet("/api/admin/messages", token);
}
