import websocket
try:
    import thread
except ImportError:
    import _thread as thread
import time

def on_message(ws, message):
    print(message)

def on_error(ws, error):
    print(error)

def on_close(ws):
    print("### closed ###")

def on_open(ws):
    def run(*args):
        for i in range(2):
            time.sleep(1)
            ws.send("Hello %d" % i)
        time.sleep(1)
        ws.close()
        print("thread terminating...")
    thread.start_new_thread(run, ())


if __name__ == "__main__":
    websocket.enableTrace(True)

    ws1 = websocket.WebSocketApp("ws://localhost:8087",
                              on_message = on_message,
                              on_error = on_error,
                              on_close = on_close)
    ws2 = websocket.WebSocketApp("ws://localhost:8087",
                                  on_message = on_message,
                                  on_error = on_error,
                                  on_close = on_close)
    ws1.on_open = on_open
    ws2.on_open = on_open
    ws1.run_forever()
    ws2.run_forever()