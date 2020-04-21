import { program } from 'commander';
import { countCommand } from './count.command';
import { pushCommand } from './push.command';
import { listCommand } from './list.command';

countCommand(program);
pushCommand(program);
listCommand(program);
program.parse(process.argv);