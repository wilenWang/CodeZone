import type { Message } from "../../lib/api";
import { MarkdownMessage } from "./MarkdownMessage";

type Props = {
  currentUserId: number;
  messages: Message[];
  failedMessages: { id: string; contentMarkdown: string }[];
  onRetry: (id: string, contentMarkdown: string) => void;
};

function formatTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

export function MessageList({ currentUserId, messages, failedMessages, onRetry }: Props) {
  const hasContent = messages.length > 0 || failedMessages.length > 0;

  return (
    <div className="message-list" role="log" aria-live="polite" aria-label="Messages">
      {!hasContent ? (
        <div className="empty-state">
          <p className="empty-state-hint">No messages yet. Send the first one.</p>
        </div>
      ) : null}
      {[...messages].reverse().map((message) => {
        const isMine = message.senderId === currentUserId;
        return (
          <article key={message.id} className={isMine ? "message mine" : "message"}>
            <div className="message-content">
              <MarkdownMessage content={message.contentMarkdown} />
            </div>
            <div className="message-meta">
              {isMine ? <span>You</span> : null}
              {isMine ? <span aria-hidden="true">·</span> : null}
              <time dateTime={message.createdAt}>{formatTime(message.createdAt)}</time>
            </div>
          </article>
        );
      })}
      {failedMessages.map((message) => (
        <article key={message.id} className="message mine failed">
          <div className="message-content">
            <MarkdownMessage content={message.contentMarkdown} />
          </div>
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
