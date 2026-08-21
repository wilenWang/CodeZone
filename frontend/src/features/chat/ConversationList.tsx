import type { Conversation, User } from "../../lib/api";
import { Avatar } from "./Avatar";

type Props = {
  conversations: Conversation[];
  user: User;
  users?: User[];
  selectedId: number | null;
  directUserId?: number | null;
  selectedDirectUserId?: number | null;
  onSelect: (id: number) => void;
  onStartDirect?: (userId: number) => void;
  onCreateGroup?: () => void;
  onOpenProfile?: () => void;
  className?: string;
  isLoading?: boolean;
  isUsersLoading?: boolean;
  error?: Error | null;
  directError?: Error | null;
  onRetry?: () => void;
};

export function ConversationList({
  conversations,
  user,
  users = [],
  selectedId,
  directUserId = null,
  selectedDirectUserId = null,
  onSelect,
  onStartDirect = () => {},
  onCreateGroup = () => {},
  onOpenProfile = () => {},
  className,
  isLoading,
  isUsersLoading,
  error,
  directError,
  onRetry,
}: Props) {
  const groupConversations = conversations.filter((conversation) => conversation.type === "group");
  const directUsers = users.filter((item) => item.id !== user.id);

  return (
    <aside className={`conversation-list ${className ?? ""}`.trim()} aria-label="Conversations">
      <button className="current-user" type="button" aria-label="个人设置" onClick={onOpenProfile}>
        <Avatar user={user} size="md" />
        <span className="current-user-details"><span className="current-user-label">当前用户</span><strong className="current-user-name">{user.displayName}</strong><span className="current-user-username">@{user.username}</span></span>
      </button>
      <button className="create-group-button" type="button" onClick={onCreateGroup}>
        新建群聊
      </button>
      <div className="list-title">群聊</div>
      {isLoading ? (
        <div className="loading-state" role="status" aria-label="Loading conversations">
          <div className="skeleton" style={{ height: 44 }} />
          <div className="skeleton" style={{ height: 44 }} />
        </div>
      ) : error ? (
        <div className="error-state">
          <p className="error-state-message">Could not load conversations</p>
          {onRetry ? <button className="primary-button" type="button" onClick={onRetry}>Retry</button> : null}
        </div>
      ) : groupConversations.length === 0 ? (
        <p className="sidebar-empty">暂无群聊</p>
      ) : (
        groupConversations.map((conversation) => (
          <button
            key={conversation.id}
            className={conversation.id === selectedId ? "conversation-row selected" : "conversation-row"}
            onClick={() => onSelect(conversation.id)}
            type="button"
            aria-current={conversation.id === selectedId ? "true" : undefined}
          >
            <span>{conversation.title ?? `Conversation ${conversation.id}`}</span>
            {conversation.unreadCount > 0 ? <span className="unread-badge" aria-label={`${conversation.unreadCount} unread`}>{conversation.unreadCount}</span> : null}
          </button>
        ))
      )}
      <div className="list-title direct-list-title">单聊</div>
      {isUsersLoading ? (
        <div className="loading-state" role="status" aria-label="Loading users"><div className="skeleton" style={{ height: 44 }} /></div>
      ) : (
        directUsers.map((item) => (
          <button
            key={item.id}
            className={
              directUserId === item.id || selectedDirectUserId === item.id
                ? "direct-user-row selected"
                : "direct-user-row"
            }
            type="button"
            onClick={() => onStartDirect(item.id)}
            disabled={directUserId !== null}
          >
            <Avatar user={item} size="sm" />
            <span className="direct-user-details"><span className="direct-user-name">{item.displayName}</span><span className="direct-user-username">@{item.username}</span></span>
          </button>
        ))
      )}
      {directError ? <p className="sidebar-error">{directError.message}</p> : null}
    </aside>
  );
}
