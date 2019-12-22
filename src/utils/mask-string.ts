/**
 * Masks a string with a given pattern.
 * @param input The input string.
 * @param indices The indices of the characters to extract from the input string.
 * @param pattern The output pattern.
 * @param replaceSymbol The symbol in the {@link pattern} to replace with the
 * extracted input characters.
 */
export function maskString(
  input: string,
  indices: number[],
  pattern: string,
  replaceSymbol: string = '#'): string {
  return indices.reduce(
    (prev, curr) => prev.replace(replaceSymbol, input.charAt(curr)),
    pattern
  );
}