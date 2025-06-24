import threading
import time
import sys


class LoadingSpinner:
    def __init__(self, text="Loading", delay=0.3):
        self.text = text
        self.delay = delay
        self.running = False
        self.thread = None
        self.spinner_chars = "|/-\\"
        self.dot_states = ["", ".", "..", "...","....","....."]

    def start(self):
        self.running = True
        self.thread = threading.Thread(target=self._spinner_task)
        self.thread.daemon = True
        self.thread.start()

    def stop(self):
        self.running = False
        if self.thread:
            self.thread.join()
        sys.stdout.write('\r' + ' ' * 80 + '\r')  # Leere Zeile
        sys.stdout.flush()

    def _spinner_task(self):
        i = 0
        j = 0
        next_time = time.perf_counter()
        while self.running:
            spinner = self.spinner_chars[i % len(self.spinner_chars)]
            dots = self.dot_states[j % len(self.dot_states)]
            output = f"\r{spinner} {self.text}{dots}"
            sys.stdout.write(output)
            sys.stdout.flush()

            i += 1
            if i % len(self.spinner_chars) == 0:
                j += 1

            # wait till next tick
            next_time += self.delay
            sleep_time = max(0, next_time - time.perf_counter())
            time.sleep(sleep_time)


# Test
if __name__ == "__main__":
    spinner = LoadingSpinner("Loading infos", delay=0.3)
    spinner.start()

    # Simulates Task
    time.sleep(5)

    spinner.stop()
    print("Finished!")

