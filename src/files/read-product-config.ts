import { readJSONSync } from 'fs-extra';
import { ProductConfig } from './product-config';
import { join, dirname } from 'path';

/**
 * Reads the product configuration file in the specified path.
 */
export function readProductConfig(path: string): ProductConfig {
  const c: ProductConfig = readJSONSync(path);
  c.dataPath = join(dirname(path), c.dataPath);
  if (c.image) {
    c.image.path = join(dirname(path), c.image.path);
  }
  return c;
}