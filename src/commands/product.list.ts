import commander from 'commander';
import { listProducts } from '../shopify/product.list';

export function listCommand(cmd: commander.Command): commander.Command {
  return cmd.command('list')
    .description('Lists all products in the store')
    .action(() => listProducts(products => console.log(products)));
}
