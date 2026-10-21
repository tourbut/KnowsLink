"""Probe the running publish fixture on the remote loopback; synthetic data only, credentials never printed."""
import json
import pathlib
import subprocess
import sys

sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent))
import remote  # noqa: E402

PROBE = r'''
set -u
U=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["url"])' "$BASE/fixture.json")
G=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["gate_pending"])' "$BASE/fixture.json")
R=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["gate_with_result"])' "$BASE/fixture.json")
code() { curl -s -o /dev/null -w '%{http_code}' "$@"; }
echo "before cleanup: $(curl -s $U/fixture/state)"
echo "retire key=present: $(curl -s -X POST "$U/fixture/retire?key=present")"
curl -s -X POST $U/fixture/seed >/dev/null
echo "retire key=present revoke=fail: $(curl -s -X POST "$U/fixture/retire?key=present&revoke=fail")"
curl -s -X POST $U/fixture/seed >/dev/null
echo "retire key=absent: $(curl -s -X POST "$U/fixture/retire?key=absent")"
for p in /privacy /terms; do
  echo "GET $p ko: $(code $U$p) h1=$(curl -s $U$p | grep -o '<h1>[^<]*</h1>')"
  echo "GET $p en: $(code -H 'Cookie: __Host-kl_lang=en' $U$p) h1=$(curl -s -H 'Cookie: __Host-kl_lang=en' $U$p | grep -o '<h1>[^<]*</h1>') korean_runs=$(curl -s -H 'Cookie: __Host-kl_lang=en' $U$p | sed 's/한국어//g' | grep -cP '[\x{AC00}-\x{D7A3}]')"
done
echo "GET / footer: $(curl -s $U/ | grep -o '<footer>.*</footer>' | grep -o 'href="/[a-z]*"' | tr '\n' ' ')"
for l in ko en; do
  S=$(curl -s -D - -o /dev/null -X POST $U/fixture/as/$l | grep -i '^set-cookie: __Host-kl_session' | sed 's/^[Ss]et-[Cc]ookie: //; s/;.*//' | tr -d '\r')
  A=$(curl -s -H "Cookie: $S" -H 'Sec-Fetch-Site: same-origin' $U/home/account)
  echo "account $l: policy_links=$(echo "$A" | grep -c 'href="/privacy"') calendar_words=$(echo "$A" | grep -ciE 'calendar|캘린더|busy|바쁜 시간')"
  P=$(curl -s -H "Cookie: $S" -H 'Sec-Fetch-Site: same-origin' $U/home/gates/$G)
  echo "pending gate $l: approve_button=$(echo "$P" | grep -c 'value="approve"') u5_text=$(echo "$P" | grep -cE '정보 공개나 일정 실행을 허용하지 않습니다|does not allow information disclosure') calendar_words=$(echo "$P" | grep -ciE 'calendar|캘린더|busy|바쁜 시간|공개되는 내용')"
  D=$(curl -s -H "Cookie: $S" -H 'Sec-Fetch-Site: same-origin' $U/home/gates/$R)
  echo "cleaned gate $l: state_denied=$(echo "$D" | grep -cE '거절|Denied') approve_button=$(echo "$D" | grep -c 'value="approve"')"
done
echo "consume cleaned gate: $(curl -s -X POST "$U/fixture/consume?gate=$R")"
echo "POST /home/calendar/connect: $(code -X POST -H "Cookie: $S" -H 'Sec-Fetch-Site: same-origin' $U/home/calendar/connect)"
echo "POST /home/calendar/disconnect: $(code -X POST -H "Cookie: $S" -H 'Sec-Fetch-Site: same-origin' $U/home/calendar/disconnect)"
echo "fixture alive: $(kill -0 $(cat $BASE/fixture.pid) && echo yes)"
'''


def main():
    sha = remote.sha_head()
    base = (remote.HERE / f'{sha}.base').read_text().strip()
    client, values = remote.connect()
    try:
        stdin, stdout, stderr = client.exec_command(f'BASE={base} bash -s', timeout=300)
        stdin.write(PROBE)
        stdin.channel.shutdown_write()
        out = stdout.read().decode('utf-8', 'replace') + stderr.read().decode('utf-8', 'replace')
        code = stdout.channel.recv_exit_status()
    finally:
        client.close()
    for value in values.values():
        if value:
            out = out.replace(value, '[redacted]')
    (remote.HERE / f'{sha}-fixture-probe.log').write_text(out + f'\nEXIT={code}\n', encoding='utf-8')
    print(out)
    return code


if __name__ == '__main__':
    raise SystemExit(main())
