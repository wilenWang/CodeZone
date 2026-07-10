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
      <div className="admin-page-inner">
        <h1>Admin</h1>
        <AdminTable title="Users" rows={users.data?.users ?? []} isLoading={users.isLoading} />
        <AdminTable
          title="Conversations"
          rows={conversations.data?.conversations ?? []}
          isLoading={conversations.isLoading}
        />
        <AdminTable
          title="Recent Messages"
          rows={messages.data?.messages ?? []}
          isLoading={messages.isLoading}
        />
      </div>
    </main>
  );
}

function AdminTable({
  title,
  rows,
  isLoading,
}: {
  title: string;
  rows: Record<string, unknown>[];
  isLoading: boolean;
}) {
  const columns = Object.keys(rows[0] ?? {});

  return (
    <section className="admin-section">
      <h2>{title}</h2>
      <div className="admin-table-wrap">
        {isLoading ? (
          <div className="loading-state" aria-label={`Loading ${title.toLowerCase()}`}>
            <div className="skeleton" style={{ height: 40 }} />
            <div className="skeleton" style={{ height: 40 }} />
            <div className="skeleton" style={{ height: 40 }} />
          </div>
        ) : rows.length === 0 ? (
          <div className="empty-state">
            <p className="empty-state-hint">No {title.toLowerCase()} found</p>
          </div>
        ) : (
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
                    <td
                      key={column}
                      className={typeof row[column] === "number" ? "tabular-nums" : ""}
                    >
                      {String(row[column] ?? "")}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </section>
  );
}
