import { shopify } from '../../shopify/shopify';
import commander from 'commander';

export function countCommand(cmd: commander.Command): commander.Command {
  return cmd.command('count')
    .description('Returns the total number of products')
    .action(() => {
      shopify.product.count().then(n => console.info(`Total number of products: ${n}`));
    });
}
