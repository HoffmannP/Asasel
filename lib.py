import subprocess

def validUser(username: str):
    user = username.split()[0]
    result = subprocess.run(['id', '-u', user], encoding='utf8', capture_output=True)
    if result.returncode != 0:
        return False, 'user does not exist'
    if int(result.stdout) < 1000:
       return False, 'can not kill this user'
    return True, user

def shell(*cmd, input=None):
    result = subprocess.run(cmd, encoding='utf8', capture_output=True, input=input)
    print(result.stdout, result.stderr)
    return (result.stdout if result.returncode == 0 else f'ERROR[{result.returncode}]: {result.stderr}').strip()
