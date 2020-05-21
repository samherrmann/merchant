import { readProducts } from '../../files/read-products';
import { progress } from '../../utils/progress';
import commander from 'commander';
import { writeProducts } from '../../files/write-products';
import { createProduct } from '../../shopify/create-product';
import { readShopConfig } from '../../files/read-shop-config';
import { readProductConfig } from '../../files/read-product-config';
import { groupArray } from '../../utils/group-array';
import { flattenArray } from '../../utils/flatten-array';

export function pushCommand(cmd: commander.Command): commander.Command {
  return cmd.command('push')
    .description('Pushes all products to the store')
    .action(async () => {
      // read the shop configuration file.
      const shopConfig = readShopConfig();
      // for each product type...
      for (const path of shopConfig.products) {
        // read the product configuration file.
        const c = readProductConfig(path);
        // get products from file that don't already exist in store.
        // each product is an array of variants.
        const products = groupArray(readProducts(c.dataPath).filter(p => !p.product_id), p => p.handle);
        // exit if no new products exist.
        if (!products.length) {
          console.info(`All "${c.type}" products already exist in the store.`);
          continue;
        }
        await progress(
          `Adding "${c.type}" products to store:`,
          products,
          async variants => {
            // add the product to the store.
            const storeProduct = await createProduct(variants, c);
            // save the store ID.
            storeProduct.variants.forEach((v, i) => {
              /* eslint-disable @typescript-eslint/camelcase */
              variants[i].product_id = storeProduct.id;
              variants[i].variant_id = v.id;
            })
          }
        );
        // write the products back to the file.
        writeProducts(c.dataPath, flattenArray(products));
      }
    });
}

