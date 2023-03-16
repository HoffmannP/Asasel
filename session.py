from fastapi import FastAPI
import subprocess

app = FastAPI()

def validUser(username: str) -> tuple[bool, str]:
    user = username.split()[0]
    result = subprocess.run(['id', '-u', user], encoding='utf8', capture_output=True)
    if result.returncode != 0:
        return False, 'user does not exist'
    if int(result.stdout) < 1000:
       return False, 'can not kill this user'
    return True, user

def shell(*cmd):
    result = subprocess.run(cmd, encoding='utf8', capture_output=True)
    return (result.stdout if result.returncode == 0 else f'ERROR: {result.stderr}').strip()

@app.get("/killall/{username}")
def whoam(username: str):
    ok, response = validUser(username)
    if not ok:
        return response
    return shell('echo', 'killall', '-u', response)


@app.get("/killsession/{username}")
def killsession(username: str):
    ok, response = validUser(username)
    if not ok:
        return response
    return shell('echo', 'pkill', '-fu', response, "'^cinnamon-session '")
