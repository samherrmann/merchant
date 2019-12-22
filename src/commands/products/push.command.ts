import { readProducts } from '../../files/read-products';
import { progress } from '../../utils/progress';
import { Command } from 'commander';
import { writeProducts } from '../../files/write-products';
import { createProduct } from '../../shopify/create-product';
import { readShopConfig } from '../../files/read-shop-config';

export function pushCommand(cmd: Command): Command {
  return cmd.command('push')
    .description('Pushes all products to the store')
    .action(async () => {

      // read the configuration file.
      const config = readShopConfig();
      // for each product type...
      for (let c of config.products) {
        // get products from file that don't already exist in store.
        const products = readProducts(c.dataPath).filter(p => !p.shopifyId)
        // exit if no new products exist.
        if (!products.length) {
          console.info(`All "${c.type}" products already exist in the store.`);
          return;
        }
        await progress(
          `Adding "${c.type}" products to store:`,
          products,
          async p => {
            // add the product to the store.
            const storeProduct = await createProduct(p, c);
            // save the store ID.
            p.shopifyId = `${storeProduct.id}`;
          }
        );
        // write the products back to the file.
        writeProducts(c.dataPath, products);
      }
    });
}

