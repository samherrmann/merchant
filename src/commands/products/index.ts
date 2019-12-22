import program from 'commander';
import { countCommand } from './count.command';
import { pushCommand } from './push.command';

countCommand(program);
pushCommand(program);
program.parse(process.argv);