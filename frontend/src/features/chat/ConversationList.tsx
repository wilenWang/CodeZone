import type { Conversation } from "../../lib/api";

type Props = {
  conversations: Conversation[];
  selectedId: number | null;
  onSelect: (id: number) => void;
  className?: string;
};

export function ConversationList({ conversations, selectedId, onSelect, className }: Props) {
  return (
    <aside className={`conversation-list ${className ?? ""}`.trim()} aria-label="Conversations">
      <div className="list-title">Chats</div>
      {conversations.map((conversation) => (
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
      ))}
    </aside>
  );
}
