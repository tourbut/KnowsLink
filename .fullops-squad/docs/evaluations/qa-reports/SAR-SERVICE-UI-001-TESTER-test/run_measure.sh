export PATH="$PATH:/c/Users/shin/AppData/Local/Programs/orca/resources/bin"
SP="$1"; JS=$(tr -d '\n' < $SP/measure.js)
: > $SP/measure2.jsonl
for w in 1440 390 320; do
 for f in $SP/fixtures/*.html; do
  n=$(basename $f .html)
  orca goto --url "file:///$f" --json >/dev/null 2>&1
  orca exec --command "set viewport $w 844" --json >/dev/null 2>&1
  out=$(orca eval --expression "$JS" --json 2>&1 | python3 -c "import sys,json;print(json.load(sys.stdin)['result']['result'])")
  echo "{\"w\":$w,\"fixture\":\"$n\",\"m\":$out}" >> $SP/measure2.jsonl
 done
done
echo done > $SP/measure2.done
