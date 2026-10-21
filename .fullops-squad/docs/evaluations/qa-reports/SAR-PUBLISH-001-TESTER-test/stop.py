"""Stop this task's fixture and list leftover publish.test processes (run via explore.py)."""
import os, json, urllib.request, subprocess
B = os.environ['BASE']; U = json.load(open(B + '/fixture.json'))['url']
print('stop', urllib.request.urlopen(urllib.request.Request(U + '/fixture/stop', method='POST', data=b'')).status)
print(subprocess.run(f'pgrep -af {B}/publish.test', shell=True, capture_output=True, text=True).stdout)
