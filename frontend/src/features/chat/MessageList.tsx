import type { Message } from "../../lib/api";
import { MarkdownMessage } from "./MarkdownMessage";

type Props = {
  currentUserId: number;
  messages: Message[];
  failedMessages: { id: string; contentMarkdown: string }[];
  onRetry: (id: string, contentMarkdown: string) => void;
};

export function MessageList({ currentUserId, messages, failedMessages, onRetry }: Props) {
  return (
    <div className="message-list">
      {[...messages].reverse().map((message) => (
        <article key={message.id} className={message.senderId === currentUserId ? "message mine" : "message"}>
          <MarkdownMessage content={message.contentMarkdown} />
        </article>
      ))}
      {failedMessages.map((message) => (
        <article key={message.id} className="message mine failed">
          <MarkdownMessage content={message.contentMarkdown} />
          <div className="message-failed-actions">
            <span>Failed to send</span>
            <button type="button" onClick={() => onRetry(message.id, message.contentMarkdown)}>
              Retry
            </button>
          </div>
        </article>
      ))}
    </div>
  );
}
