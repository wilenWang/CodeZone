import type { Conversation } from "../../lib/api";

type Props = {
  conversations: Conversation[];
  selectedId: number | null;
  onSelect: (id: number) => void;
  className?: string;
  isLoading?: boolean;
  error?: Error | null;
  onRetry?: () => void;
};

export function ConversationList({
  conversations,
  selectedId,
  onSelect,
  className,
  isLoading,
  error,
  onRetry,
}: Props) {
  return (
    <aside className={`conversation-list ${className ?? ""}`.trim()} aria-label="Conversations">
      <div className="list-title">Chats</div>
      {isLoading ? (
        <div className="loading-state" role="status" aria-label="Loading conversations">
          <div className="skeleton" style={{ height: 44 }} />
          <div className="skeleton" style={{ height: 44 }} />
          <div className="skeleton" style={{ height: 44 }} />
        </div>
      ) : error ? (
        <div className="error-state">
          <p className="error-state-message">Could not load conversations</p>
          {onRetry ? (
            <button className="primary-button" type="button" onClick={onRetry}>
              Retry
            </button>
          ) : null}
        </div>
      ) : conversations.length === 0 ? (
        <div className="empty-state">
          <p className="empty-state-hint">No conversations yet</p>
        </div>
      ) : (
        conversations.map((conversation) => (
          <button
            key={conversation.id}
            className={conversation.id === selectedId ? "conversation-row selected" : "conversation-row"}
            onClick={() => onSelect(conversation.id)}
            type="button"
            aria-current={conversation.id === selectedId ? "true" : undefined}
          >
            <span>{conversation.title ?? `Conversation ${conversation.id}`}</span>
            {conversation.unreadCount > 0 ? (
              <span className="unread-badge" aria-label={`${conversation.unreadCount} unread`}>
                {conversation.unreadCount}
              </span>
            ) : null}
          </button>
        ))
      )}
    </aside>
  );
}
