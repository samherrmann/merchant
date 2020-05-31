/**
 * Groups the items of the provided array by a given key. If the `groupKey`
 * function returns `undefined` for an item, then that item is returned in its
 * own group. The order of the items is maintained.
 */
export function groupArray<T>(arr: T[], groupKey: (obj: T) => string | undefined): T[][] {
  const groups = arr.reduce((acc, curr, index) => {
    const key = groupKey(curr);
    if (key) {
      const group = acc.get(key) || [];
      group.push(curr)
      acc.set(key, group);
    } else {
      acc.set(`${index}`, [curr]);
    }
    return acc;
  }, new Map<string, T[]>());
  return Array.from(groups.values());
}