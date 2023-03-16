#!/usr/bin/env python3

import importlib
import logging
import subprocess

import fastapi
import fastapi.middleware.cors
import uvicorn  # type: ignore

import session

app = fastapi.FastAPI()
app.add_middleware(
    fastapi.middleware.cors.CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['POST', 'GET'],
    allow_headers=['Content-Type'],
)

@app.get("/whoami")
def whoam():
    result = subprocess.run('whoami', encoding='utf8', capture_output=True)
    return (result.stdout if result.returncode == 0 else result.stderr).strip()

app.mount("/session", session.app)

if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s | %(levelname)s:     %(message)s', level=logging.INFO)
    logging.info('Logging started')
    uvicorn.run(app, host='0.0.0.0')
