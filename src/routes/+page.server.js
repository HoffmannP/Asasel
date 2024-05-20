import child_process from 'child_process'
import { commander, commands } from '../lib/commands.js'

function parameterlist (command) {
    return Array.from(new Set(
        command.filter(
            c => (c[0] === '$')
        ).map(
            c => c.substring(1)
        )
    ))
}

function parameterlistDict(commandlist) {
    return Object.fromEntries(
        Object.entries(commandlist).map(
            ([ k, cmd ]) => [ k, parameterlist(cmd) ]
        )
    )
}

function composeCommand(formdata) {
    const command = commands[formdata.get('command')]
    const parameters = Object.fromEntries(parameterlist(command).map(
        (param, i) => [param, formdata.get(`param${i}`)]
    ))
    console.debug(command, parameters)
    return [
        commander,
        command.map(arg => arg[0] === '$'
            ? parameters[arg.substring(1)]
            : arg
        )
    ]
}

export async function load (event) {
    return { commands: parameterlistDict(commands) }
}

export const actions = {
    default: async function ({ request }) {
        const formdata = await request.formData()
        const command = composeCommand(formdata).flat(Infinity)
        console.debug(command)
        const returnobject = child_process.execSync(command.join(" "))
        const stdout = returnobject.toString().trim()
        console.debug(stdout)
        if (stdout != '') {
            return stdout
        }
        return null
    }
}