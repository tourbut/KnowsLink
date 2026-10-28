"""Compare server observations before/during/after the 0.1.1 relay switch; writes redacted preservation.json (counts, digests, booleans only). Raw observations stay in .devtools/autosetup-001-ops/obs-*.json."""
import json, pathlib
src = pathlib.Path('D:/workspace/KnowsLink/.devtools/autosetup-001-ops')
root = pathlib.Path(__file__).parent
STAGES = ['baseline', 'before', 'frozen', 'after', 'pre-override', 'after-override', 'final', 'final2']
o = {n: json.loads((src / f'obs-{n}.json').read_text()) for n in STAGES}
A = []
def check(name, ok, detail=None):
    A.append({'assertion': name, 'pass': bool(ok), **({'detail': detail} if detail is not None else {})})
base, b, f, a, po, ao, z, z2 = (o[n] for n in STAGES)
VOL = {'Rates'}  # request-time records; changes with every request, not a permission change
for k in b['digests']:
    if k not in VOL:
        check(f'state map {k} digest equal at all {len(STAGES)} stages', len({o[n]['digests'][k] for n in o}) == 1)
check('state map counts equal at all stages except volatile Rates', all({k: v for k, v in o[n]['counts'].items() if k not in VOL} == {k: v for k, v in base['counts'].items() if k not in VOL} for n in o), {k: v for k, v in base['counts'].items() if k in ('Members', 'Owners', 'Agents', 'Friends', 'Pairs', 'Messages')})
check('Remote digests/counts, inbox and subscription metadata equal at all stages', all(o[n][k] == base[k] for n in o for k in ['remote_inbox_count', 'subscription_metadata', 'remote_digests', 'remote_counts']))
check('frozen equals before for every state digest incl. Rates (writer stopped)', b['digests'] == f['digests'] and b['remote_digests'] == f['remote_digests'])
check('schema sha256 equal at all stages', len({o[n]['schema_sha256'] for n in o}) == 1)
check('HumanModel=1 and goose2/guard1/rows1 at all stages (no migration)', all(o[n]['human_model'] == 1 and json.loads(o[n]['goose_trigger']) == {'goose': 2, 'guard': 1, 'rows': 1} for n in o))
fe = {n: o[n]['files'] for n in o}
check('.env, tunnel.uuid, access.aud and tunnel/config.yml hash+mode unchanged at all stages (env and ingress untouched)', all(fe[n] == fe['baseline'] for n in o) and fe['baseline']['.env']['mode'] == '0o600')
check('tunnel credential json hash+mode equal at all stages', all(o[n]['private_tunnel_files'] == base['private_tunnel_files'] for n in o))
check('state and backups directory modes 0700 at all stages', all(o[n]['state_mode'] == '0o700' and o[n]['backups_mode'] == '0o700' for n in o))
bi = base['backup_inventory']
new = [k for k in a['backup_inventory'] if k not in bi]
check('existing 7 backups preserved (hash/mode/bytes); exactly one new backup, unchanged afterwards', len(bi) == 7 and len(new) == 1 and all(a['backup_inventory'][k] == v for k, v in bi.items()) and all(o[n]['backup_inventory'] == a['backup_inventory'] for n in ['after', 'pre-override', 'after-override', 'final', 'final2']), {'before_count': len(bi), 'after_count': len(a['backup_inventory'])})
check('new backup is mode 0600', a['backup_inventory'][new[0]]['mode'] == '0o600', {'bytes': a['backup_inventory'][new[0]]['bytes'], 'sha256': a['backup_inventory'][new[0]]['sha256']})
cmap = lambda x: {c['name']: c for c in x['containers']}
c0 = cmap(base)
sig = lambda c: {k: c[k] for k in ('image', 'id', 'running', 'health', 'started', 'restarts')}
check('before: every container identical to baseline', all(sig(cmap(b)[k]) == sig(c0[k]) for k in c0) and set(cmap(b)) == set(c0))
check('frozen: only relay stopped (writer stop); every other container identical to baseline', not cmap(f)['/knowslink-relay-1']['running'] and all(sig(cmap(f)[k]) == sig(c0[k]) for k in c0 if k != '/knowslink-relay-1') and set(cmap(f)) == set(c0))
for n in ('after', 'pre-override', 'after-override', 'final', 'final2'):
    changed = sorted(k for k in c0 if sig(cmap(o[n])[k]) != sig(c0[k]))
    check(f'{n}: only relay and migrate containers changed', changed == ['/knowslink-migrate-1', '/knowslink-relay-1'], changed)
check('Postgres and KnowsLink cloudflared containers unchanged through final2', all(sig(cmap(z2)[k]) == sig(c0[k]) for k in ('/knowslink-postgres-1', '/knowslink-cloudflared-1')))
check('only new container since baseline is the retained restore fixture', sorted(k for k in cmap(z2) if k not in c0) == ['/knowslink-autosetup011-rehearse-20261028'])
check('relay switch: relay image changed; relay healthy with restarts 0 at final2', c0['/knowslink-relay-1']['image'] != cmap(z2)['/knowslink-relay-1']['image'] and cmap(z2)['/knowslink-relay-1']['health'] == 'healthy' and cmap(z2)['/knowslink-relay-1']['restarts'] == 0)
check('override removal did not recreate relay (container id and start time same before/after)', sig(cmap(po)['/knowslink-relay-1']) == sig(cmap(ao)['/knowslink-relay-1']) == sig(cmap(z2)['/knowslink-relay-1']))
check('shared orca/s8/mcp public statuses equal at all stages', all(o[n]['shared'] == base['shared'] for n in o))
pids = {n: o[n]['host_cloudflared_pids'] for n in o}
check('host cloudflared (cloudflared-orca user service, not KnowsLink) PID: same through after-override; new at final (service restarted 2026-10-28T13:21:38Z after its config.yml was edited at 13:21:07Z by someone else)', all(pids[n] == pids['baseline'] for n in STAGES[:6]) and pids['final'] == pids['final2'] and pids['final'] != pids['baseline'], pids)
check('listeners 127.0.0.1:8080 and 5432 restored; frozen had only 5432', base['listeners'] == a['listeners'] == z2['listeners'] and f['listeners'] == ['127.0.0.1:5432'])
check('no other active DB writer at frozen and later', all(json.loads(o[n]['db_connections'])['active_other_writers'] == 0 for n in STAGES[2:]))
OLD, NEW = 'bf67dfcbb4f7164149901d847edc0aa0e5783dc7', 'c6dfc39ac4f6dabddc5864ad611437b84ded383e'
check('checkout bf67dfcb at baseline; c6dfc39a and clean from before on', base['head'] == OLD and not base['dirty'] and all(o[n]['head'] == NEW and not o[n]['dirty'] for n in STAGES[1:]))
epochs = [o[n]['epoch'] for n in STAGES]
check('epoch never decreased', epochs == sorted(epochs), epochs)
out = {'utc_stages': {n: o[n]['utc'] for n in STAGES}, 'passed': sum(x['pass'] for x in A), 'total': len(A), 'assertions': A}
(root / 'preservation.json').write_text(json.dumps(out, indent=2, ensure_ascii=False) + '\n', encoding='utf-8')
print(out['passed'], '/', out['total'])
for x in A:
    if not x['pass']:
        print('FAIL', x['assertion'], x.get('detail'))
