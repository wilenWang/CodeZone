import { describe, expect, it } from "vitest";
import { authHeader } from "./api";

describe("authHeader", () => {
  it("returns bearer header when token exists", () => {
    expect(authHeader("abc")).toEqual({ Authorization: "Bearer abc" });
  });

  it("returns empty object without token", () => {
    expect(authHeader(null)).toEqual({});
  });
});
