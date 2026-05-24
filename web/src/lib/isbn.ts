/**
 * ISBN-13 チェックサムアルゴリズム (GS1): 奇数桁×1, 偶数桁×3 の合計が 10 の倍数であれば valid
 */
export function isValidIsbn13(isbn: string): boolean {
  const digits = isbn.replace(/-/g, "");
  if (!/^\d{13}$/.test(digits)) return false;
  const sum = digits.split("").reduce((acc, d, i) => acc + Number(d) * (i % 2 === 0 ? 1 : 3), 0);
  return sum % 10 === 0;
}
