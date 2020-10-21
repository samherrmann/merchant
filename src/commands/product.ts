import { program } from 'commander';
import { countCommand } from './product.count';
import { pushCommand } from './product.push';
import { listCommand } from './product.list';

countCommand(program);
pushCommand(program);
listCommand(program);
program.parse(process.argv);