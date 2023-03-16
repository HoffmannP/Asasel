from fastapi import FastAPI
import lib

app = FastAPI()

def endSession(username):
    return ['killall', '-u', username, 'cinnamon-session']

@app.get('/killall/{username}')
def whoam(username: str):
    ok, response = lib.validUser(username)
    if not ok:
        return response
    return lib.shell('killall', '-u', response)

@app.get('/killsession/{username}')
def killsession(username: str):
    ok, response = lib.validUser(username)
    if not ok:
        return response
    return lib.shell(*endSession(response))

def changePassword(username: str, password: str):
    return lib.shell('passwd', username, input=f'{password}\n{password}\n')

@app.get('/simplepasswd/{username}')
def simplepasswd(username: str):
    ok, response = lib.validUser(username)
    if not ok:
        return response
    return changePassword(response, response)

@app.get('/hardpasswd/{username}')
def hardpasswd(username: str):
    ok, response = lib.validUser(username)
    if not ok:
        return response
    return changePassword(response, f'{response.capitalize()}{response.capitalize()}')

@app.get('/setlimit/{username}/{time}')
def setlimit(username: str, time: str):
    ok, response = lib.validUser(username)
    if not ok:
        return response

    if time[0] == '+':
        time = f'now {time}'

    return lib.shell('at', time, input=' '.join(endSession(response)))

@app.get('/clearlimit/{username}')
def clearlimit(username: str):
    ok, response = lib.validUser(username)
    if not ok:
        return response

    match = ' '.join(endSession(username))

    for line in lib.shell('atq').split('\n'):
        atid = line.split('\t', 1)[0]
        command = lib.shell('at', '-c', atid).rsplit('\n', 1)[1]
        if command == match:
            lib.shell('atrm', atid)
            return command
        else:
            print(command, match)
    return 'No session limit found'

@app.get('/shutdown')
def shutdown():
    return lib.shell('shutdown', '-P')