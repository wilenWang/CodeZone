import { useState } from "react";

type Props = {
  disabled: boolean;
  onSend: (text: string) => void;
};

export function MessageComposer({ disabled, onSend }: Props) {
  const [value, setValue] = useState("");

  return (
    <form
      className="composer"
      onSubmit={(event) => {
        event.preventDefault();
        const trimmed = value.trim();
        if (!trimmed) {
          return;
        }
        onSend(trimmed);
        setValue("");
      }}
    >
      <textarea
        aria-label="Message"
        value={value}
        disabled={disabled}
        onChange={(event) => setValue(event.target.value)}
      />
      <button type="submit" disabled={disabled}>
        Send
      </button>
    </form>
  );
}
