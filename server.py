#!/usr/bin/env python3

import datetime
import logging
import subprocess

import fastapi
import fastapi.middleware.cors
import fastapi.responses
import jwt
import pam
import uvicorn  # type: ignore

import session

app = fastapi.FastAPI()
app.add_middleware(
    fastapi.middleware.cors.CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['POST', 'GET'],
    allow_headers=['Content-Type'],
)

SECRET = 'ChangeThisSecretToYourOwn'
COOKIE_NAME = 'asasel_jwt'
JWT_ALGORITHM = 'HS256'

def validJwt(request):
    if 'authorization' in request.headers and request.headers('authorization').startsWith('Bearer '):
        token = request.headers('authorization')[7:]
    elif COOKIE_NAME in request.cookies:
        token = request.cookies[COOKIE_NAME]
    else:
        return False
    try:
        payload = jwt.decode(token, key=SECRET, algorithms=JWT_ALGORITHM, options={
            'verify_signature': True,
            'require': ['iss', 'aud', 'exp', 'iat', 'nbf'],
            'verify_iss': True,
            'verify_aud': True,
            'verify_exp': True,
            'verify_iat': True,
            'verify_nbf': True},
            audience='Asasel',
            issuer='Asasel',
            leeway=60)
        print(payload['payload'])
        return True
    except jwt.exceptions as excep:
        print(excep)
        return False

@app.middleware('http')
async def checkAuthBearerToken(
        request: fastapi.Request,
        call_next):

    if not (request.url.path.startswith('/login/') or validJwt(request)):
        return fastapi.responses.RedirectResponse('/login/', status_code=401)
    return await call_next(request)


@app.get('/login/')
def loginForm():
    loginForm = '''
<form method="POST">
    <input type="text" name="username" placeholder="username"><br>
    <input type="password" name="password"  placeholder="password"><br>
    <button>Login</button>
</form>
    '''
    return fastapi.Response(content=loginForm, media_type="text/html")

@app.post('/login/')
def sendToken(response: fastapi.Response, username: str = fastapi.Form(), password: str = fastapi.Form()):
    if pam.authenticate(username, password):
        now = datetime.datetime.now()
        timestamp = lambda dt: int(dt.timestamp())
        payload = {'iss': 'Asasel',
            'sub': 'Asasel',
            'aud': 'Asasel',
            'exp': timestamp(now + datetime.timedelta(weeks=6)),
            'iat': timestamp(now),
            'nbf': timestamp(now),
            'payload': {'username': username}}
        token = jwt.encode(payload, SECRET, algorithm=JWT_ALGORITHM)
        response.set_cookie(key=COOKIE_NAME, value=token)
        return token
    return fastapi.responses.RedirectResponse('/login/', status_code=401)

@app.get('/whoami')
def whoam():
    result = subprocess.run('whoami', encoding='utf8', capture_output=True)
    return (result.stdout if result.returncode == 0 else result.stderr).strip()

app.mount('/session', session.app)

@app.get('/{catchAll:path}')
def catchAll(catchAll: str):
    response = f'Path "{catchAll}" does not exist'
    print(response)
    return response


if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s | %(levelname)s:     %(message)s', level=logging.INFO)
    logging.info('Logging started')
    uvicorn.run(app, host='0.0.0.0')
