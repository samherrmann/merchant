import { readJSONSync } from 'fs-extra';
import { ShopConfig } from './shop-config';

/**
 * Reads the shopctl.json file from the current working directory.
 */
export function readShopConfig(): ShopConfig {
  return readJSONSync('./shopctl.json')
}