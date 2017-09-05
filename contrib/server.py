import socket, struct

connection = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
connection.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
connection.bind(('0.0.0.0', 8000))
connection.listen(10)
while True:
    conn, _ = connection.accept()
    while True:
        hdr = conn.recv(2)
        if not hdr:
            break
        (l,) = struct.unpack(">H", hdr)
        data = conn.recv(l)
        if not data:
            break
        conn.send(hdr + data)
