import { Command } from 'commander';
import { listProducts } from '../../shopify/list-products';

export function listCommand(cmd: Command): Command {
  return cmd.command('list')
    .description('Lists all products in the store')
    .action(() => listProducts(products => console.log(products)));
}
