export function flattenArray<T>(arr: T[][]): T[] {
  return ([] as T[]).concat.apply([], arr);
}