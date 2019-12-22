import { shopify } from '../../shopify/shopify';
import { Command } from 'commander';

export function countCommand(cmd: Command): Command {
  return cmd.command('count')
    .description('Returns the total number of products')
    .action(() => {
      shopify.product.count().then(n => console.info(`Total number of products: ${n}`));
    });
}
