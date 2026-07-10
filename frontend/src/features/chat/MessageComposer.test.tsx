import { fireEvent, render } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { MessageComposer } from "./MessageComposer";

describe("MessageComposer", () => {
  it("sends message on submit", () => {
    const onSend = vi.fn();
    const { container } = render(<MessageComposer disabled={false} onSend={onSend} />);
    const input = container.querySelector("textarea");
    expect(input).not.toBeNull();
    fireEvent.change(input!, { target: { value: "Hello" } });
    const button = container.querySelector("button");
    fireEvent.click(button!);
    expect(onSend).toHaveBeenCalledWith("Hello");
  });

  it("sends message on Enter", () => {
    const onSend = vi.fn();
    const { container } = render(<MessageComposer disabled={false} onSend={onSend} />);
    const input = container.querySelector("textarea");
    expect(input).not.toBeNull();
    fireEvent.change(input!, { target: { value: "Hello" } });
    fireEvent.keyDown(input!, { key: "Enter", code: "Enter", shiftKey: false });
    expect(onSend).toHaveBeenCalledWith("Hello");
  });

  it("does not send on Shift+Enter", () => {
    const onSend = vi.fn();
    const { container } = render(<MessageComposer disabled={false} onSend={onSend} />);
    const input = container.querySelector("textarea");
    expect(input).not.toBeNull();
    fireEvent.change(input!, { target: { value: "Line 1" } });
    fireEvent.keyDown(input!, { key: "Enter", code: "Enter", shiftKey: true });
    expect(onSend).not.toHaveBeenCalled();
  });

  it("does not send empty messages", () => {
    const onSend = vi.fn();
    const { container } = render(<MessageComposer disabled={false} onSend={onSend} />);
    const button = container.querySelector("button");
    fireEvent.click(button!);
    expect(onSend).not.toHaveBeenCalled();
  });

  it("disables send while disabled", () => {
    const { container } = render(<MessageComposer disabled onSend={() => {}} />);
    const button = container.querySelector("button") as HTMLButtonElement;
    const textarea = container.querySelector("textarea") as HTMLTextAreaElement;
    expect(button.disabled).toBe(true);
    expect(textarea.disabled).toBe(true);
  });
});
