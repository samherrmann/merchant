import { readFileSync } from 'fs-extra';
import parseSync from 'csv-parse/lib/sync';
import { Product } from './product';
import { csvDelimiter } from './csv-delimiter';

/**
 * Reads the products from the provided `path`.
 */
export function readProducts(path: string): Product[] {
  const csv = readFileSync(path).toString();
  return parseSync(csv, {
    delimiter: csvDelimiter,
    columns: true,
    skip_empty_lines: true
  });
}
