import { QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { LoginPage } from "../features/auth/LoginPage";
import { ChatPage } from "../features/chat/ChatPage";
import type { User } from "../lib/api";
import { queryClient } from "../lib/query";

export function App() {
  const [token, setToken] = useState(() => localStorage.getItem("chat.token"));
  const [user, setUser] = useState<User | null>(null);

  return (
    <QueryClientProvider client={queryClient}>
      {token && user ? (
        <ChatPage token={token} user={user} />
      ) : (
        <LoginPage
          onLogin={(nextToken, nextUser) => {
            localStorage.setItem("chat.token", nextToken);
            setToken(nextToken);
            setUser(nextUser);
          }}
        />
      )}
    </QueryClientProvider>
  );
}
