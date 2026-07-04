import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { MessageComposer } from "./MessageComposer";

describe("MessageComposer", () => {
  it("submits trimmed message text", () => {
    const onSend = vi.fn();
    render(<MessageComposer disabled={false} onSend={onSend} />);
    fireEvent.change(screen.getByRole("textbox"), { target: { value: " hello " } });
    fireEvent.click(screen.getByRole("button", { name: /send/i }));
    expect(onSend).toHaveBeenCalledWith("hello");
  });
});
