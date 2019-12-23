import { readFileSync } from 'fs-extra';
import { dirname, join } from 'path';
import parseSync from 'csv-parse/lib/sync';
import { Product } from './product';
import { csvDelimiter } from './csv-delimiter';

/**
 * Reads the products from the provided `path`. The path may be provided
 * relative to another file/directory by specifying the `relativeTo` argument.
 */
export function readProducts(path: string, relativeTo?: string): Product[] {
  if (relativeTo) {
    path = join(dirname(relativeTo), path);
  }
  const csv = readFileSync(path).toString();
  return parseSync(csv, {
    delimiter: csvDelimiter,
    columns: true,
    skip_empty_lines: true
  });
}
