import { fireEvent, render } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { Conversation, User } from "../../lib/api";
import { ConversationList } from "./ConversationList";

describe("ConversationList", () => {
  const conversations: Conversation[] = [
    { id: 1, type: "group", title: "Project Room", lastMessageId: null, lastMessageAt: null, unreadCount: 3 },
    { id: 2, type: "direct", title: "Direct with Bob", lastMessageId: null, lastMessageAt: null, unreadCount: 0 },
  ];
  const user: User = { id: 1, username: "alice", displayName: "Alice", avatarUrl: null, userType: "human" };
  const users: User[] = [user, { id: 2, username: "bob", displayName: "Bob", avatarUrl: null, userType: "human" }];

  function renderList(overrides: Partial<React.ComponentProps<typeof ConversationList>> = {}) {
    return render(
      <ConversationList
        conversations={conversations}
        user={user}
        users={users}
        selectedId={null}
        onSelect={() => {}}
        {...overrides}
      />,
    );
  }

  it("renders only group conversations in the group section", () => {
    const { container } = renderList();
    const rows = container.querySelectorAll(".conversation-row");
    expect(rows).toHaveLength(1);
    expect(rows[0]?.textContent).toContain("Project Room");
    expect(container.textContent).not.toContain("Direct with Bob");
  });

  it("shows unread badges for group conversations", () => {
    const { container } = renderList();
    const badges = container.querySelectorAll(".unread-badge");
    expect(badges).toHaveLength(1);
    expect(badges[0]?.textContent).toBe("3");
  });

  it("calls onSelect when a group conversation is clicked", () => {
    const onSelect = vi.fn();
    const { container } = renderList({ onSelect });
    fireEvent.click(container.querySelector(".conversation-row")!);
    expect(onSelect).toHaveBeenCalledWith(1);
  });

  it("lists other users as direct-chat targets and starts a direct chat", () => {
    const onStartDirect = vi.fn();
    const { container } = renderList({ onStartDirect });
    expect(container.textContent).toContain("Bob");
    expect(container.querySelectorAll(".direct-user-row")).toHaveLength(1);
    fireEvent.click(container.querySelector(".direct-user-row")!);
    expect(onStartDirect).toHaveBeenCalledWith(2);
  });

  it("shows an empty group state", () => {
    const { container } = renderList({ conversations: [] });
    expect(container.textContent).toContain("暂无群聊");
  });

  it("shows error state with retry", () => {
    const onRetry = vi.fn();
    const { container } = renderList({ error: new Error("fail"), onRetry });
    const retryButton = [...container.querySelectorAll("button")].find((button) => button.textContent === "Retry");
    fireEvent.click(retryButton!);
    expect(onRetry).toHaveBeenCalled();
  });
});
