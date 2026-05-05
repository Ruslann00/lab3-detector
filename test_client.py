import time
import requests

PPROF_URL = "http://localhost:6060/debug/pprof/"

def main():
    print("Checking pprof endpoint...")

    while True:
        try:
            response = requests.get(PPROF_URL, timeout=2)
            print("Status:", response.status_code)
        except Exception as error:
            print("Error:", error)

        time.sleep(1)

if __name__ == "__main__":
    main()