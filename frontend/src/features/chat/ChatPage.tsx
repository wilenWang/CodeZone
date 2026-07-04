import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import {
  listConversations,
  listMessages,
  sendMessage,
  type Conversation,
  type Message,
  type User,
} from "../../lib/api";
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
  });
  const activeConversation = (conversations.data?.conversations ?? []).find(
    (item: Conversation) => item.id === activeId,
  );

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
            <header className="chat-header">{activeConversation?.title ?? `Conversation ${activeId}`}</header>
            <MessageList currentUserId={user.id} messages={messages.data?.messages ?? []} />
            <MessageComposer disabled={send.isPending} onSend={(text) => send.mutate(text)} />
          </>
        ) : (
          <div className="empty-state">Select a conversation</div>
        )}
      </section>
    </main>
  );
}
