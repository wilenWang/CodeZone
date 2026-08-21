import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { cleanup, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import type { User } from "../../lib/api";
import { ChatPage } from "./ChatPage";

const user: User = { id: 1, username: "alice", displayName: "Alice", avatarUrl: null, userType: "human" };

const mocks = vi.hoisted(() => ({
  listConversations: vi.fn(),
  listUsers: vi.fn(),
  ensureDirectConversation: vi.fn(),
  listMessages: vi.fn(),
  markRead: vi.fn(),
  sendMessage: vi.fn(),
  connectEvents: vi.fn(),
}));

vi.mock("../../lib/api", async () => {
  const actual = await vi.importActual<typeof import("../../lib/api")>("../../lib/api");
  return {
    ...actual,
    listConversations: mocks.listConversations,
    listUsers: mocks.listUsers,
    ensureDirectConversation: mocks.ensureDirectConversation,
    listMessages: mocks.listMessages,
    markRead: mocks.markRead,
    sendMessage: mocks.sendMessage,
  };
});

vi.mock("../../lib/ws", () => ({
  connectEvents: mocks.connectEvents,
}));

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <ChatPage token="test-token" user={user} onUserChange={() => {}} />
    </QueryClientProvider>,
  );
}

describe("ChatPage", () => {
  afterEach(() => {
    cleanup();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    mocks.connectEvents.mockReturnValue(() => {});
    mocks.listConversations.mockResolvedValue({
      conversations: [
        {
          id: 1,
          type: "group",
          title: "Project Room",
          lastMessageId: null,
          lastMessageAt: null,
          unreadCount: 2,
        },
      ],
    });
    mocks.listMessages.mockResolvedValue({ messages: [] });
    mocks.markRead.mockResolvedValue({ ok: true });
    mocks.listUsers.mockResolvedValue({
      users: [user, { id: 2, username: "bob", displayName: "Bob", avatarUrl: null, userType: "human" }],
    });
    mocks.ensureDirectConversation.mockResolvedValue({
      id: 2,
      type: "direct",
      title: null,
      lastMessageId: null,
      lastMessageAt: null,
      unreadCount: 0,
    });
  });

  it("highlights the direct user whose conversation is active", async () => {
    const { container } = renderPage();
    await screen.findByText("Bob");
    const bobRow = container.querySelector<HTMLButtonElement>(".direct-user-row");
    expect(bobRow).toBeTruthy();
    bobRow!.click();
    // wait until the direct-conversation mutation has fully settled (row re-enabled)
    await waitFor(() => expect(bobRow!.disabled).toBe(false));
    expect(bobRow!.className).toContain("selected");
  });

  it("shows the current user's display name and username", async () => {
    renderPage();
    expect(screen.getByText("Alice")).toBeTruthy();
    expect(screen.getByText("@alice")).toBeTruthy();
  });

  it("marks an active conversation read when it has unread messages", async () => {
    renderPage();
    await waitFor(() => expect(mocks.markRead).toHaveBeenCalledWith("test-token", 1));
  });

  it("does not mark read when the active conversation has no unread messages", async () => {
    mocks.listConversations.mockResolvedValue({
      conversations: [
        {
          id: 1,
          type: "group",
          title: "Project Room",
          lastMessageId: null,
          lastMessageAt: null,
          unreadCount: 0,
        },
      ],
    });
    renderPage();
    await waitFor(() => expect(mocks.listConversations).toHaveBeenCalled());
    expect(mocks.markRead).not.toHaveBeenCalled();
  });
});
