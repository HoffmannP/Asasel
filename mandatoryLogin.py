import json
from typing import Any, Callable

from starlette.types import ASGIApp, Scope, Receive, Send
from starlette.requests import Request
from starlette.responses import Response, RedirectResponse

LOGIN_FORM = '''<form method="POST">
    <input type="text" name="username" placeholder="username"><br>
    <input type="password" name="password"  placeholder="password"><br>
    <button> Login </button>
</form>'''

class MandatoryLoginMiddleware:
    def __init__(self,
                 app: ASGIApp,
                 authFkt: Callable[[str, str], bool],
                 tokenHandler: Any,
                 location: str = '/login',
                 form: str = LOGIN_FORM):
        self.app = app
        self.authFkt = authFkt
        self.location = location
        self.form = form
        self.token = tokenHandler

    def redirectToLogin(self) -> Response:
        return RedirectResponse(url=self.location)

    def sendHtml(self, body: str) -> Response:
        return Response(body, media_type='text/html')

    def sendJson(self, obj: Any) -> Response:
        return Response(json.dumps(obj), media_type='text/html')

    async def __call__(self, scope: Scope, receive: Receive, send: Send) -> None:
        if scope["type"] != "http":
            response = self.app
        else:
            request = Request(scope, receive)
            if request.url.path == self.location:
                if request.method == 'GET':
                    response = self.sendHtml(self.form)
                elif request.method == 'POST':
                    form = await request.form()
                    if type(form['username']) is str and type(form['password']) is str and \
                        self.authFkt(form['username'], form['password']):
                        response = self.token.assign(form['username'])
                    else:
                        response = self.redirectToLogin()
                else:
                    response = self.redirectToLogin()
            else:
                if self.token.validate(request):
                    response = self.app
                else:
                    response = self.redirectToLogin()
        return await response(scope, receive, send)
