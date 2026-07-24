import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import {
  listConversations,
  listMessages,
  markRead,
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
  const [mobileShowChat, setMobileShowChat] = useState(false);
  const activeIdRef = useRef<number | null>(null);

  const conversations = useQuery({
    queryKey: ["conversations"],
    queryFn: () => listConversations(token),
  });

  const activeId = selectedId ?? conversations.data?.conversations[0]?.id ?? null;
  activeIdRef.current = activeId;

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

  const markReadMutation = useMutation({
    mutationFn: (conversationId: number) => markRead(token, conversationId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["conversations"] });
    },
  });

  useEffect(() => {
    if (activeId && activeConversation && activeConversation.unreadCount > 0) {
      markReadMutation.mutate(activeId);
    }
  }, [activeId, activeConversation?.unreadCount, markReadMutation]);

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
          const id = activeIdRef.current;
          if (id) {
            void queryClient.invalidateQueries({ queryKey: ["messages", id] });
          }
        }
      },
    );
  }, [token, queryClient]);

  function retryFailedMessage(id: string, contentMarkdown: string) {
    setFailedMessages((items) => items.filter((item) => item.id !== id));
    send.mutate(contentMarkdown);
  }

  function handleSelect(id: number) {
    setSelectedId(id);
    setMobileShowChat(true);
  }

  function handleBackToList() {
    setMobileShowChat(false);
  }

  const listHiddenClass = mobileShowChat ? "mobile-hidden" : "";
  const chatVisibleClass = mobileShowChat ? "mobile-visible" : "";

  return (
    <main className="chat-layout">
      <ConversationList
        className={listHiddenClass}
        conversations={conversations.data?.conversations ?? []}
        error={conversations.error}
        isLoading={conversations.isLoading}
        onRetry={() => conversations.refetch()}
        selectedId={activeId}
        onSelect={handleSelect}
      />
      <section className={`chat-panel ${chatVisibleClass}`.trim()}>
        {activeId ? (
          <>
            <header className="chat-header">
              <div className="chat-header-start">
                <button
                  className="back-button"
                  type="button"
                  onClick={handleBackToList}
                  aria-label="Back to conversations"
                >
                  ←
                </button>
                <span className="chat-header-title">
                  {activeConversation?.title ?? `Conversation ${activeId}`}
                </span>
              </div>
              <span
                className={connected ? "ws-status connected" : "ws-status"}
                role="status"
                aria-live="polite"
              >
                <span className="ws-status-dot" aria-hidden="true" />
                {connected ? "Live" : "Reconnecting"}
              </span>
            </header>
            {messages.isLoading ? (
              <div className="loading-state" role="status" aria-label="Loading messages">
                <div className="skeleton" style={{ height: 64 }} />
                <div className="skeleton" style={{ height: 64, width: "80%" }} />
                <div className="skeleton" style={{ height: 64, width: "60%" }} />
              </div>
            ) : messages.error ? (
              <div className="error-state">
                <p className="error-state-message">Could not load messages</p>
                <button className="primary-button" type="button" onClick={() => messages.refetch()}>
                  Retry
                </button>
              </div>
            ) : (
              <MessageList
                currentUserId={user.id}
                failedMessages={failedMessages}
                messages={messages.data?.messages ?? []}
                onRetry={retryFailedMessage}
              />
            )}
            <MessageComposer disabled={send.isPending} onSend={(text) => send.mutate(text)} />
          </>
        ) : (
          <div className="empty-state">
            <div>
              <p className="empty-state-title">Select a conversation</p>
              <p className="empty-state-hint">Choose a chat from the list to start messaging</p>
            </div>
          </div>
        )}
      </section>
    </main>
  );
}
