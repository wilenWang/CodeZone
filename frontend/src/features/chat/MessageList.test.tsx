import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { Message } from "../../lib/api";
import { MessageList } from "./MessageList";

describe("MessageList", () => {
  const messages: Message[] = [
    {
      id: 1,
      conversationId: 1,
      senderId: 1,
      contentMarkdown: "Hello",
      contentPlain: "Hello",
      createdAt: "2026-07-10T10:00:00Z",
    },
    {
      id: 2,
      conversationId: 1,
      senderId: 2,
      contentMarkdown: "Hi there",
      contentPlain: "Hi there",
      createdAt: "2026-07-10T10:01:00Z",
    },
  ];

  it("renders messages", () => {
    const { container } = render(
      <MessageList
        currentUserId={1}
        messages={messages}
        failedMessages={[]}
        onRetry={() => {}}
      />,
    );
    const list = screen.getByRole("log", { name: "Messages" });
    expect(list).toBeDefined();
    const articles = screen.getAllByRole("article");
    expect(articles.length).toBe(2);
    expect(screen.getByText("Hello")).toBeDefined();
    expect(screen.getByText("Hi there")).toBeDefined();
    const time = container.querySelector('time[datetime="2026-07-10T10:00:00Z"]');
    expect(time).not.toBeNull();
    expect(time?.tagName).toBe("TIME");
    expect(screen.getByText("You")).toBeDefined();
  });

  it("renders failed messages with retry button", () => {
    const onRetry = vi.fn();
    render(
      <MessageList
        currentUserId={1}
        messages={[]}
        failedMessages={[{ id: "failed-1", contentMarkdown: "Oops" }]}
        onRetry={onRetry}
      />,
    );
    expect(screen.getByText("Oops")).toBeDefined();
    expect(screen.getByText("Failed to send")).toBeDefined();
    const retryButton = screen.getByRole("button", { name: "Retry" });
    fireEvent.click(retryButton);
    expect(onRetry).toHaveBeenCalledWith("failed-1", "Oops");
  });
});
