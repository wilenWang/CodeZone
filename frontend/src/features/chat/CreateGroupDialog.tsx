import { useState } from "react";
import type { User } from "../../lib/api";
import { Avatar } from "./Avatar";

type Props = {
  users: User[];
  isCreating: boolean;
  error: Error | null;
  onClose: () => void;
  onCreate: (title: string, memberIds: number[]) => void;
};

export function CreateGroupDialog({ users, isCreating, error, onClose, onCreate }: Props) {
  const [title, setTitle] = useState("");
  const [memberIds, setMemberIds] = useState<number[]>([]);
  const [validationError, setValidationError] = useState<string | null>(null);

  function toggleMember(userId: number) {
    setMemberIds((ids) => (ids.includes(userId) ? ids.filter((id) => id !== userId) : [...ids, userId]));
    setValidationError(null);
  }

  function submit(event: React.FormEvent) {
    event.preventDefault();
    const trimmedTitle = title.trim();
    if (!trimmedTitle) {
      setValidationError("请输入群聊名称");
      return;
    }
    if (memberIds.length < 2) {
      setValidationError("请至少选择两位成员");
      return;
    }
    onCreate(trimmedTitle, memberIds);
  }

  return (
    <div className="dialog-backdrop" role="presentation">
      <section className="group-dialog" role="dialog" aria-modal="true" aria-labelledby="group-dialog-title">
        <div className="group-dialog-header">
          <h2 id="group-dialog-title">新建群聊</h2>
          <button className="dialog-close" type="button" onClick={onClose} aria-label="关闭">
            ×
          </button>
        </div>
        <form onSubmit={submit}>
          <label className="field-label">
            群聊名称
            <input
              className="text-input"
              value={title}
              onChange={(event) => {
                setTitle(event.target.value);
                setValidationError(null);
              }}
              disabled={isCreating}
              autoFocus
            />
          </label>
          <fieldset className="group-members" disabled={isCreating}>
            <legend>选择成员（至少两位）</legend>
            {users.map((user) => (
              <label key={user.id} className="group-member-option">
                <input
                  type="checkbox"
                  checked={memberIds.includes(user.id)}
                  onChange={() => toggleMember(user.id)}
                />
                <Avatar user={user} size="sm" />
                <span>
                  <strong>{user.displayName}</strong>
                  <small>@{user.username}</small>
                </span>
              </label>
            ))}
          </fieldset>
          {validationError ? <p className="error-text">{validationError}</p> : null}
          {error ? <p className="error-text">{error.message}</p> : null}
          <div className="dialog-actions">
            <button className="secondary-button" type="button" onClick={onClose} disabled={isCreating}>
              取消
            </button>
            <button className="primary-button" type="submit" disabled={isCreating}>
              {isCreating ? "创建中..." : "创建群聊"}
            </button>
          </div>
        </form>
      </section>
    </div>
  );
}
