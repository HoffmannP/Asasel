#!env python3

import logging
import fastapi
import fastapi.staticfiles
import uvicorn  # type: ignore
import server

base = fastapi.FastAPI()
base.mount('/api/v1', server.app, name='api')
base.mount('/', fastapi.staticfiles.StaticFiles(directory='static', html=True), name='static')

if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s | %(levelname)s:     %(message)s', level=logging.INFO)
    logging.info('Logging started')
    uvicorn.run(base, host='0.0.0.0', port=80)
