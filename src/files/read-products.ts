import { readFileSync } from 'fs-extra';
import parseSync from 'csv-parse/lib/sync';
import { Product } from './product';
import { csvDelimiter } from './csv-delimiter';
import { groupArray } from '../utils/group-array';
import { ProductConfig } from './product-config';
import mustache from 'mustache';

/**
 * Reads the products from the provided `path`.
 */
export function readProducts(c: ProductConfig, newOnly?: boolean): Product[][] {
  const csv = readFileSync(c.dataPath).toString();
  let products: Product[] = parseSync(csv, {
    delimiter: csvDelimiter,
    columns: true,
    // eslint-disable-next-line @typescript-eslint/camelcase
    skip_empty_lines: true
  });
  if (newOnly) {
    // eslint-disable-next-line @typescript-eslint/camelcase
    products = products.filter(p => !p.product_id);
  }
  return groupArray(products, p => mustache.render(c.title, p));
}
