import { describe, expect, it } from "vitest";
import { isValidIsbn13 } from "./isbn";

describe("isValidIsbn13", () => {
	it("valid ISBN-13 を受け入れる", () => {
		expect(isValidIsbn13("9784798142470")).toBe(true);
	});

	it("ハイフン区切りを受け入れる", () => {
		expect(isValidIsbn13("978-4-7981-4247-0")).toBe(true);
	});

	it("チェックディジットが不正な場合は false を返す", () => {
		expect(isValidIsbn13("9784798142471")).toBe(false);
	});

	it("数字以外を含む場合は false を返す", () => {
		expect(isValidIsbn13("978479814247X")).toBe(false);
	});

	it("桁数が 13 未満の場合は false を返す", () => {
		expect(isValidIsbn13("978479814247")).toBe(false);
	});
});
