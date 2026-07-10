import { fireEvent, render } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { Conversation } from "../../lib/api";
import { ConversationList } from "./ConversationList";

describe("ConversationList", () => {
  const conversations: Conversation[] = [
    { id: 1, type: "group", title: "Project Room", lastMessageId: null, lastMessageAt: null, unreadCount: 3 },
    { id: 2, type: "direct", title: "Direct with Bob", lastMessageId: null, lastMessageAt: null, unreadCount: 0 },
  ];

  it("renders conversation titles", () => {
    const { container } = render(
      <ConversationList conversations={conversations} selectedId={null} onSelect={() => {}} />,
    );
    const rows = container.querySelectorAll(".conversation-row");
    expect(rows[0]?.textContent).toContain("Project Room");
    expect(rows[1]?.textContent).toContain("Direct with Bob");
  });

  it("shows unread badge only for conversations with unread messages", () => {
    const { container } = render(
      <ConversationList conversations={conversations} selectedId={null} onSelect={() => {}} />,
    );
    const badges = container.querySelectorAll(".unread-badge");
    expect(badges.length).toBe(1);
    expect(badges[0]?.textContent).toBe("3");
  });

  it("calls onSelect when a conversation is clicked", () => {
    const onSelect = vi.fn();
    const { container } = render(
      <ConversationList conversations={conversations} selectedId={null} onSelect={onSelect} />,
    );
    const firstRow = container.querySelector(".conversation-row");
    expect(firstRow).not.toBeNull();
    fireEvent.click(firstRow!);
    expect(onSelect).toHaveBeenCalledWith(1);
  });

  it("marks selected conversation", () => {
    const { container } = render(
      <ConversationList conversations={conversations} selectedId={2} onSelect={() => {}} />,
    );
    const selected = container.querySelector(".selected");
    expect(selected).not.toBeNull();
    expect(selected?.textContent).toContain("Direct with Bob");
  });

  it("shows loading skeletons", () => {
    const { container } = render(
      <ConversationList conversations={[]} selectedId={null} onSelect={() => {}} isLoading />,
    );
    expect(container.querySelectorAll(".skeleton").length).toBeGreaterThan(0);
  });

  it("shows empty state when there are no conversations", () => {
    const { container } = render(
      <ConversationList conversations={[]} selectedId={null} onSelect={() => {}} />,
    );
    expect(container.textContent).toContain("No conversations yet");
  });

  it("shows error state with retry", () => {
    const onRetry = vi.fn();
    const { container } = render(
      <ConversationList
        conversations={[]}
        selectedId={null}
        onSelect={() => {}}
        error={new Error("fail")}
        onRetry={onRetry}
      />,
    );
    expect(container.textContent).toContain("Could not load conversations");
    const retryButton = container.querySelector("button");
    fireEvent.click(retryButton!);
    expect(onRetry).toHaveBeenCalled();
  });
});
