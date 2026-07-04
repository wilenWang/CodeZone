import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import {
  listConversations,
  listMessages,
  sendMessage,
  type Conversation,
  type Message,
  type User,
} from "../../lib/api";
import { connectEvents } from "../../lib/ws";
import { ConversationList } from "./ConversationList";
import { MessageComposer } from "./MessageComposer";
import { MessageList } from "./MessageList";

type Props = {
  token: string;
  user: User;
};

export function ChatPage({ token, user }: Props) {
  const queryClient = useQueryClient();
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [connected, setConnected] = useState(false);
  const [failedMessages, setFailedMessages] = useState<{ id: string; contentMarkdown: string }[]>([]);
  const conversations = useQuery({
    queryKey: ["conversations"],
    queryFn: () => listConversations(token),
  });
  const activeId = selectedId ?? conversations.data?.conversations[0]?.id ?? null;
  const messages = useQuery({
    queryKey: ["messages", activeId],
    queryFn: () => listMessages(token, activeId!),
    enabled: activeId !== null,
  });
  const send = useMutation({
    mutationFn: (text: string) => sendMessage(token, activeId!, text),
    onSuccess: (message: Message) => {
      queryClient.setQueryData<{ messages: Message[] }>(["messages", activeId], (old) => ({
        messages: [message, ...(old?.messages ?? [])],
      }));
      void queryClient.invalidateQueries({ queryKey: ["conversations"] });
    },
    onError: (_error, text) => {
      setFailedMessages((items) => [...items, { id: crypto.randomUUID(), contentMarkdown: text }]);
    },
  });
  const activeConversation = (conversations.data?.conversations ?? []).find(
    (item: Conversation) => item.id === activeId,
  );

  useEffect(() => {
    return connectEvents(
      token,
      (event) => {
        if (event.type === "message.created") {
          void queryClient.invalidateQueries({ queryKey: ["messages"] });
        }
        if (event.type === "conversation.updated") {
          void queryClient.invalidateQueries({ queryKey: ["conversations"] });
        }
      },
      (nextConnected) => {
        setConnected(nextConnected);
        if (nextConnected) {
          void queryClient.invalidateQueries({ queryKey: ["conversations"] });
          if (activeId) {
            void queryClient.invalidateQueries({ queryKey: ["messages", activeId] });
          }
        }
      },
    );
  }, [token, queryClient, activeId]);

  function retryFailedMessage(id: string, contentMarkdown: string) {
    setFailedMessages((items) => items.filter((item) => item.id !== id));
    send.mutate(contentMarkdown);
  }

  return (
    <main className="chat-layout">
      <ConversationList
        conversations={conversations.data?.conversations ?? []}
        selectedId={activeId}
        onSelect={setSelectedId}
      />
      <section className="chat-panel">
        {activeId ? (
          <>
            <header className="chat-header">
              <span>{activeConversation?.title ?? `Conversation ${activeId}`}</span>
              <span className={connected ? "ws-status connected" : "ws-status"}>
                {connected ? "Live" : "Reconnecting"}
              </span>
            </header>
            <MessageList
              currentUserId={user.id}
              failedMessages={failedMessages}
              messages={messages.data?.messages ?? []}
              onRetry={retryFailedMessage}
            />
            <MessageComposer disabled={send.isPending} onSend={(text) => send.mutate(text)} />
          </>
        ) : (
          <div className="empty-state">Select a conversation</div>
        )}
      </section>
    </main>
  );
}
