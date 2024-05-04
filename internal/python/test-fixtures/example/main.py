import requests
from example import message


def handler(event, ctx):
    message('bob')
    resp = requests.get('https://example.com')
    print('Got request status code', resp.status_code)


if __name__ == '__main__':
    message('joe')
