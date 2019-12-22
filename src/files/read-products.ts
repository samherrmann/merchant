import { readFileSync } from 'fs-extra';
import parseSync from 'csv-parse/lib/sync';
import { Product } from './product';
import { csvDelimiter } from './csv-delimiter';

export function readProducts(path: string): Product[] {
  const csv = readFileSync(path).toString();
  return parseSync(csv, {
    delimiter: csvDelimiter,
    columns: true,
    skip_empty_lines: true
  });
}

