import { describe, expect, it } from "vitest";
import { wsUrl } from "./ws";

describe("wsUrl", () => {
  it("converts http origin to ws url", () => {
    expect(wsUrl("http://localhost:5173", "/api/ws")).toBe("ws://localhost:5173/api/ws");
  });

  it("converts https origin to wss url", () => {
    expect(wsUrl("https://example.com", "/api/ws")).toBe("wss://example.com/api/ws");
  });
});
