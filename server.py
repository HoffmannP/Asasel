#!/usr/bin/env python3

import logging
import subprocess

import fastapi
import fastapi.middleware.cors
import fastapi.responses
import uvicorn  # type: ignore

import session

app = fastapi.FastAPI()
app.add_middleware(
    fastapi.middleware.cors.CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['POST', 'GET'],
    allow_headers=['Content-Type'],
)

TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQXNhc2VsIn0.FNVcvEC8GqY86L4TxgEQErHsqEdRsQU3aub4BAZf_0Q'

@app.middleware('http')
async def checkAuthBearerToken(
        request: fastapi.Request,
        call_next):

    if not request.url.path.startswith('/login/') and (
            request.headers.get('authorization') != f'Bearer {TOKEN}' and request.cookies.get('asasel_auth') != TOKEN):
        return fastapi.responses.RedirectResponse('/login/', status_code=307)
    return await call_next(request)

@app.get('/login/{username}')
def token(response: fastapi.Response, username: str, password: str):
    if username == password:
        response.set_cookie(key='asasel_auth', value=TOKEN)
        return TOKEN
    return ''

@app.get('/whoami')
def whoam():
    result = subprocess.run('whoami', encoding='utf8', capture_output=True)
    return (result.stdout if result.returncode == 0 else result.stderr).strip()

app.mount('/session', session.app)

if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s | %(levelname)s:     %(message)s', level=logging.INFO)
    logging.info('Logging started')
    uvicorn.run(app, host='0.0.0.0')
