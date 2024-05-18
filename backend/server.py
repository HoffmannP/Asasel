#!/usr/bin/env python3

import typing
import logging
import fastapi
import fastapi.middleware.cors
import fastapi_utils.tasks  # type: ignore
import uvicorn  # type: ignore

PERIOD = 5 * 60

app = fastapi.FastAPI()
app.add_middleware(
    fastapi.middleware.cors.CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['POST', 'GET'],
    allow_headers=['Content-Type'],
)


@app.get('/a')
async def a() -> typing.Dict[str, str]:
    return { 'error': 'All ticket data is to be send via POST requests'}

@app.get('/b/{year}/{week}')
async def exams_get(year: int, week: int) -> typing.Dict[str, str]:
    return { 'calender_week': f'{week}.{year}' }

@app.post('/c')
async def c(request: fastapi.Request) -> bool:
    body = await request.body()
    return len(body) > 100

@app.on_event('startup')
@fastapi_utils.tasks.repeat_every(seconds=PERIOD)
async def regular_check() -> None:
    logging.info('check')


if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s | %(levelname)s:     %(message)s', level=logging.INFO)
    logging.info('Logging started')
    uvicorn.run(app, host='0.0.0.0')
