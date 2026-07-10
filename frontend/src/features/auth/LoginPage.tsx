import { useState } from "react";
import { devLogin, type User } from "../../lib/api";

const seedUsers = ["alice", "bob", "carol"];

const avatarColors: Record<string, string> = {
  alice: "#e4e4e7",
  bob: "#d4d4d8",
  carol: "#e4e4e7",
};

type Props = {
  onLogin: (token: string, user: User) => void;
};

export function LoginPage({ onLogin }: Props) {
  const [username, setUsername] = useState("alice");
  const [error, setError] = useState<string | null>(null);
  const [inputError, setInputError] = useState(false);
  const [loadingSeed, setLoadingSeed] = useState<string | null>(null);

  async function loginAs(name: string) {
    setError(null);
    setInputError(false);
    setLoadingSeed(name);
    try {
      const result = await devLogin(name);
      onLogin(result.token, result.user);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoadingSeed(null);
    }
  }

  return (
    <main className="login-page">
      <section className="login-panel">
        <h1>Vibework Chat</h1>
        <div className="seed-list" role="list">
          {seedUsers.map((name) => (
            <button
              key={name}
              className="seed-button"
              disabled={loadingSeed !== null}
              onClick={() => void loginAs(name)}
              role="listitem"
              type="button"
            >
              <span
                className="seed-avatar"
                style={{ background: avatarColors[name] ?? "#e4e4e7" }}
                aria-hidden="true"
              >
                {name[0]?.toUpperCase()}
              </span>
              <span>Continue as {name}</span>
              {loadingSeed === name ? <span className="spinner" aria-hidden="true" /> : null}
            </button>
          ))}
        </div>
        <form
          onSubmit={(event) => {
            event.preventDefault();
            const trimmed = username.trim();
            if (!trimmed) {
              setInputError(true);
              return;
            }
            void loginAs(trimmed);
          }}
        >
          <label className="field-label">
            Username
            <input
              className={inputError ? "text-input text-input--error" : "text-input"}
              value={username}
              onChange={(event) => {
                setUsername(event.target.value);
                setInputError(false);
              }}
              aria-invalid={inputError}
              aria-describedby={inputError ? "username-error" : undefined}
            />
            {inputError ? (
              <p id="username-error" className="error-text">
                Please enter a username
              </p>
            ) : null}
          </label>
          <button className="primary-button" type="submit" disabled={loadingSeed !== null}>
            {loadingSeed === username ? <span className="spinner" aria-hidden="true" /> : null}
            Login
          </button>
        </form>
        {error ? <div className="error-banner">{error}</div> : null}
      </section>
    </main>
  );
}
