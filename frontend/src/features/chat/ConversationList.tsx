import type { Conversation } from "../../lib/api";

type Props = {
  conversations: Conversation[];
  selectedId: number | null;
  onSelect: (id: number) => void;
};

export function ConversationList({ conversations, selectedId, onSelect }: Props) {
  return (
    <aside className="conversation-list">
      <div className="list-title">Chats</div>
      {conversations.map((conversation) => (
        <button
          key={conversation.id}
          className={conversation.id === selectedId ? "conversation-row selected" : "conversation-row"}
          onClick={() => onSelect(conversation.id)}
        >
          <span>{conversation.title ?? `Conversation ${conversation.id}`}</span>
          {conversation.unreadCount > 0 ? <strong>{conversation.unreadCount}</strong> : null}
        </button>
      ))}
    </aside>
  );
}
