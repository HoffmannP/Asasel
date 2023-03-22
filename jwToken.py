from datetime import timedelta, datetime
from typing import cast, Union, Optional

from starlette.requests import Request
from starlette.responses import Response
import jwt


class JwToken:
    def __init__(self,
                 secret: str,
                 algo: str = 'HS256',
                 cookie_name: str ='jwt',
                 subject: str ='jwt handler',
                 iss: Optional[str] = None,
                 sub: Optional[str] = None,
                 aud: Optional[str] = None,
                 lifespan: Union[int, timedelta] = timedelta(weeks=2)) -> None:
        self.secret = secret
        self.algo = algo
        self.cookie_name = cookie_name
        self.iss = subject if iss is None else iss
        self.sub = subject if sub is None else sub
        self.aud = subject if aud is None else aud
        self.lifespan = lifespan if type(lifespan) is timedelta else timedelta(weeks=cast(int, lifespan))

    def assign(self, username: str) -> Response:
        now = datetime.now()
        def timestamp(dt): return int(dt.timestamp())
        payload = {'iss': self.iss,
                'sub': self.sub,
                'aud': self.aud,
                'exp': timestamp(now + self.lifespan),
                'iat': timestamp(now),
                'nbf': timestamp(now),
                'payload': {'username': username}}
        token = jwt.encode(payload, self.secret, algorithm=self.algo)
        response = Response(f'"{token}"', media_type='application/json')
        response.set_cookie(self.cookie_name, token)
        return response

    def validate(self, request: Request) -> bool:
        if 'authorization' in request.headers and request.headers['authorization'].startswith('Bearer '):
            token = request.headers['authorization'][7:]
        elif self.cookie_name in request.cookies:
            token = request.cookies[self.cookie_name]
        else:
            return False
        try:
            jwt.decode(token, key=self.secret, algorithms=[self.algo], options={
                'verify_signature': True,
                'require': ['iss', 'aud', 'exp', 'iat', 'nbf'],
                'verify_iss': True,
                'verify_aud': True,
                'verify_exp': True,
                'verify_iat': True,
                'verify_nbf': True},
            audience=self.aud,
            issuer=self.iss,
            leeway=60)
            return True
        except jwt.exceptions as excep:  # type: ignore
            print(excep)
            return False
