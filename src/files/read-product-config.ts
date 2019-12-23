import { readJSONSync } from 'fs-extra';
import { ProductConfig } from './product-config';

/**
 * Reads the product configuration file in the specified path.
 */
export function readProductConfig(path: string): ProductConfig {
  return readJSONSync(path);
}