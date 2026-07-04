import type { Message } from "../../lib/api";
import { MarkdownMessage } from "./MarkdownMessage";

type Props = {
  currentUserId: number;
  messages: Message[];
};

export function MessageList({ currentUserId, messages }: Props) {
  return (
    <div className="message-list">
      {[...messages].reverse().map((message) => (
        <article key={message.id} className={message.senderId === currentUserId ? "message mine" : "message"}>
          <MarkdownMessage content={message.contentMarkdown} />
        </article>
      ))}
    </div>
  );
}
