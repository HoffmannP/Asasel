from fastapi import FastAPI
import lib
import datetime

app = FastAPI()

def endSession(username):
    return ['killall', '-u', username, 'cinnamon-session']

def genHardPass(username):
    return f'{username.capitalize()}{username.capitalize()}'

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
    return changePassword(response, genHardPass(username))

@app.get('/setlimit/{username}/{time}')
def setlimit(username: str, time: str):
    ok, response = lib.validUser(username)
    if not ok:
        return response

    if time[0] == '+':
        time = f'now {time}'

    return lib.shell('at', time, input=' '.join(endSession(response)))

@app.get('/getlimit/{username}')
def getlimit(username: str):
    ok, response = lib.validUser(username)
    if not ok:
        return response
    match = ' '.join(endSession(username))

    for line in lib.shell('atq').split('\n'):
        atid = line.split('\t', 1)[0]
        if len(line) == 0:
            continue
        command = lib.shell('at', '-c', atid).rsplit('\n', 1)[1]
        if command == match:
            then = datetime.datetime.strptime(line.rsplit('\t', 1)[1].rsplit(' ', 2)[0], '%c')
            diff = int((then - datetime.datetime.now()).total_seconds())
            return f'{int(diff / 3600):02d}:{int(diff / 60) % 60:02d}:{diff % 60:02d}'
        else:
            print(command, match)
    return 'No session limit found'

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
    return lib.shell('shutdown', '-P', 'now')
