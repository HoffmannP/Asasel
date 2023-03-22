#!/usr/bin/env python3

from datetime import timedelta
import logging
import subprocess

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse

import pam  # type: ignore
import uvicorn  # type: ignore

from jwToken import JwToken
from mandatoryLogin import MandatoryLoginMiddleware
from session import app as sessionApp

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['POST', 'GET'],
    allow_headers=['Content-Type'],
)

app.add_middleware(
    MandatoryLoginMiddleware,
    authFkt=pam.authenticate,
    tokenHandler=JwToken(
        secret='ChangeThisSecretToYourOwn',
        cookie_name='asasel_jwt',
        subject='Asasel web remote control',
        lifespan=timedelta(weeks=6)
    )
)

@app.get('/favicon.ico')
def favicon():
    return FileResponse('favicon.ico')

@app.get('/whoami')
def whoam():
    result = subprocess.run('whoami', encoding='utf8', capture_output=True)
    return (result.stdout if result.returncode == 0 else result.stderr).strip()


app.mount('/session', sessionApp)

@app.get('/{catchAll:path}')
def catchAll(catchAll: str):
    response = f'Path "{catchAll}" does not exist'
    print(response)
    return response

def main():
    logging.basicConfig(format='%(asctime)s | %(levelname)s:     %(message)s', level=logging.INFO)
    logging.info('Logging started')
    uvicorn.run(app, host='0.0.0.0', port=2727)


if __name__ == '__main__':
    main()
