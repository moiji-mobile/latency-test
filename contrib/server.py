import socket, time

connection = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
connection.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
connection.bind(('0.0.0.0', 8000))
connection.listen(10)
while True:
    conn, _ = connection.accept()
    while True:
        data = conn.recv(16)
        if not data:
            break
        conn.send(data)
