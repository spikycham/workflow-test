import { describe, it, expect } from "vitest";
import { add } from "./math";

describe("add", () => {
    it("should return the value plus one", () => {
        expect(add(1)).toBe(3);
    });
});
