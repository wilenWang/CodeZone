import { useState } from "react";
import { devLogin, type User } from "../../lib/api";

const seedUsers = ["alice", "bob", "carol"];

type Props = {
  onLogin: (token: string, user: User) => void;
};

export function LoginPage({ onLogin }: Props) {
  const [username, setUsername] = useState("alice");
  const [error, setError] = useState<string | null>(null);

  async function loginAs(name: string) {
    setError(null);
    try {
      const result = await devLogin(name);
      onLogin(result.token, result.user);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    }
  }

  return (
    <main className="login-page">
      <section className="login-panel">
        <h1>Vibework Chat</h1>
        <div className="seed-list">
          {seedUsers.map((name) => (
            <button key={name} onClick={() => void loginAs(name)}>
              Continue as {name}
            </button>
          ))}
        </div>
        <form
          onSubmit={(event) => {
            event.preventDefault();
            void loginAs(username);
          }}
        >
          <input value={username} onChange={(event) => setUsername(event.target.value)} />
          <button type="submit">Login</button>
        </form>
        {error ? <p className="error-text">{error}</p> : null}
      </section>
    </main>
  );
}
