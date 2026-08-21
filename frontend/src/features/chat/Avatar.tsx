import { useState } from "react";
import type { User } from "../../lib/api";

type Props = {
  user: Pick<User, "username" | "displayName" | "avatarUrl">;
  size?: "sm" | "md" | "lg";
};

function colorFor(username: string) {
  let hash = 0;
  for (const char of username) hash = (hash * 31 + char.charCodeAt(0)) | 0;
  return `hsl(${Math.abs(hash) % 360} 32% 58%)`;
}

export function Avatar({ user, size = "md" }: Props) {
  const [failed, setFailed] = useState(false);
  const initial = user.displayName.trim().charAt(0).toUpperCase() || user.username.charAt(0).toUpperCase();
  if (user.avatarUrl && !failed) {
    return <img className={`avatar avatar-${size}`} src={user.avatarUrl} alt={user.displayName} onError={() => setFailed(true)} />;
  }
  return <span className={`avatar avatar-${size}`} style={{ backgroundColor: colorFor(user.username) }} aria-label={`${user.displayName} avatar`}>{initial}</span>;
}
