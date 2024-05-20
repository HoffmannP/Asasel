export const commander = 'src/lib/commands.sh'

const forbidden_execs = [
    '/usr/bin/firefox',
    '/usr/bin/google-chrome',
    '/usr/bin/minecraft-launcher',
]
const username = 'linus'
const safe_password = 'Linus1Linus'

export const commands = {
    'Einloggen erlauben': [ 'chpwd', username, username ],
    'Einloggen verbieten': [ 'chpwd', username, safe_password ],
    'Arbeitsmodus': [ 'chmod', 'o-x', forbidden_execs ],
    'Spielmodus': [ 'chmod', 'o+x', forbidden_execs ],
    'Zeitlimit': [ 'setTimelimit', username, '$Dauer' ],
    'Zeitlimit anzeigen': [ 'showTimelimit' ],
    'Zeitlimit löschen': [ 'rmTimelimit' ],
}