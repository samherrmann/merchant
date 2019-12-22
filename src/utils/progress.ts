import { Presets, Bar } from 'cli-progress';

/**
 * Executes a task for each item in a list and displays the progress in the console.
 * @param label Label describing the task being processed.
 * @param taskList List of items to process.
 * @param forEach Function executed for each item in the task list.
 */
export async function progress<T>(
  label: string,
  taskList: T[],
  forEach: (item: T) => void | Promise<void>
): Promise<void> {

  console.info(label);
  const bar = new Bar({}, Presets.rect);
  bar.start(taskList.length, 0);
  for (const [index, item] of taskList.entries()) {
    await forEach(item);
    bar.update(index + 1);
  }
  bar.stop();
}