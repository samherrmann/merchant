import commander from 'commander';
import { listCollections } from '../shopify/smart-collection.list';


export function listCommand(cmd: commander.Command): commander.Command {
  return cmd.command('list')
    .description('Lists all collections in the store')
    .action(() => listCollections().then(collections => console.log(collections)));
}
