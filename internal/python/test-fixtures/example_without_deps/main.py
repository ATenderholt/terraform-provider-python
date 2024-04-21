from example import message

def handler(event, ctx):
    message('bob')


if __name__ == '__main__':
    message('joe')
