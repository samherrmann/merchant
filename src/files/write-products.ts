import { Product } from './product';
import stringifySync from 'csv-stringify/lib/sync';
import { writeFileSync } from 'fs-extra';
import { csvDelimiter } from './csv-delimiter';

export function writeProducts(path: string, products: Product[]): void {
  const csv = stringifySync(products, {
    delimiter: csvDelimiter,
    header: true
  })
  writeFileSync(path, csv);
}