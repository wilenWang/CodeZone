import { useQuery } from "@tanstack/react-query";
import { adminConversations, adminMessages, adminUsers } from "../../lib/api";

type Props = {
  token: string;
};

export function AdminPage({ token }: Props) {
  const users = useQuery({ queryKey: ["admin", "users"], queryFn: () => adminUsers(token) });
  const conversations = useQuery({
    queryKey: ["admin", "conversations"],
    queryFn: () => adminConversations(token),
  });
  const messages = useQuery({ queryKey: ["admin", "messages"], queryFn: () => adminMessages(token) });

  return (
    <main className="admin-page">
      <h1>Admin</h1>
      <AdminTable title="Users" rows={users.data?.users ?? []} />
      <AdminTable title="Conversations" rows={conversations.data?.conversations ?? []} />
      <AdminTable title="Recent Messages" rows={messages.data?.messages ?? []} />
    </main>
  );
}

function AdminTable({ title, rows }: { title: string; rows: Record<string, unknown>[] }) {
  const columns = Object.keys(rows[0] ?? {});

  return (
    <section className="admin-section">
      <h2>{title}</h2>
      <table>
        <thead>
          <tr>
            {columns.map((column) => (
              <th key={column}>{column}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={index}>
              {columns.map((column) => (
                <td key={column}>{String(row[column] ?? "")}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
