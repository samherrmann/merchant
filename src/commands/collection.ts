import { program } from 'commander';
import { listCommand } from './collection.list';

listCommand(program);
program.parse(process.argv);
