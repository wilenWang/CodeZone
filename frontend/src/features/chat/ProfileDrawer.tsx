import { useEffect, useState } from "react";
import type { User } from "../../lib/api";
import { Avatar } from "./Avatar";

type Props = {
  user: User;
  saving: boolean;
  error: Error | null;
  onClose: () => void;
  onSaveName: (displayName: string) => void;
  onUploadAvatar: (file: File) => void;
};

export function ProfileDrawer({ user, saving, error, onClose, onSaveName, onUploadAvatar }: Props) {
  const [displayName, setDisplayName] = useState(user.displayName);
  const [preview, setPreview] = useState<string | null>(null);
  const [fileError, setFileError] = useState<string | null>(null);
  useEffect(() => () => { if (preview) URL.revokeObjectURL(preview); }, [preview]);
  const previewUser = { ...user, avatarUrl: preview ?? user.avatarUrl };
  function choose(file: File | undefined) {
    if (!file) return;
    if (!["image/jpeg", "image/png", "image/webp"].includes(file.type) || file.size > 5 * 1024 * 1024) { setFileError("请选择不超过 5 MB 的 JPEG、PNG 或 WebP 图片"); return; }
    if (preview) URL.revokeObjectURL(preview);
    setFileError(null); setPreview(URL.createObjectURL(file)); onUploadAvatar(file);
  }
  return <div className="profile-drawer-backdrop" onMouseDown={onClose}>
    <aside className="profile-drawer" aria-label="个人设置" onMouseDown={(event) => event.stopPropagation()}>
      <div className="profile-drawer-header"><div><p>账户</p><h2>个人设置</h2></div><button className="dialog-close" type="button" onClick={onClose} aria-label="关闭">×</button></div>
      <div className="profile-avatar-editor"><Avatar user={previewUser} size="lg" /><label className="secondary-button">更换头像<input className="sr-only" type="file" accept="image/jpeg,image/png,image/webp" onChange={(event) => choose(event.target.files?.[0])} disabled={saving} /></label></div>
      {fileError ? <p className="error-text">{fileError}</p> : null}
      <label className="field-label">显示名<input className="text-input" value={displayName} onChange={(event) => setDisplayName(event.target.value)} disabled={saving} /></label>
      <div className="profile-username"><span>用户名</span><strong>@{user.username}</strong></div>
      {error ? <p className="error-text">{error.message}</p> : null}
      <button className="primary-button" type="button" disabled={saving || !displayName.trim()} onClick={() => onSaveName(displayName.trim())}>{saving ? "保存中..." : "保存资料"}</button>
    </aside>
  </div>;
}
