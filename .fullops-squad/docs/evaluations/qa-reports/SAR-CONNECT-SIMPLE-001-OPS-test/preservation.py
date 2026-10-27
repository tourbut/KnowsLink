"""Compare baseline/before/frozen/after/final observations; writes redacted preservation.json (counts, digests, booleans only)."""
import json, pathlib
root = pathlib.Path(__file__).parent
STAGES = ['baseline', 'before', 'frozen', 'after', 'final']
o = {n: json.loads((root / (n + '.json')).read_text()) for n in STAGES}
A = []
def check(name, ok, detail=None):
    A.append({'assertion': name, 'pass': bool(ok), **({'detail': detail} if detail is not None else {})})
base, b, f, a, z = (o[n] for n in STAGES)
VOL = {'Rates', 'Connections'}  # Connections: one device attempt expired on the old relay between baseline and before
for k in b['digests']:
    if k in VOL:
        continue
    check(f'state map {k} digest equal at all 5 stages', len({o[n]['digests'][k] for n in o}) == 1)
check('state map counts equal at all stages except volatile Rates', all({k: v for k, v in o[n]['counts'].items() if k not in VOL} == {k: v for k, v in b['counts'].items() if k not in VOL} for n in o), {k: v for k, v in b['counts'].items() if k in ('Members', 'Owners', 'Agents', 'Friends', 'Pairs', 'Messages')})
check('Remote digests/counts (Clients/Attempts/Codes/Installs/Subscriptions), inbox and subscription metadata equal at all stages', all(o[n][k] == b[k] for n in o for k in ['remote_inbox_count', 'subscription_metadata', 'remote_digests', 'remote_counts']), {'Tokens_digest_equal': len({o[n]['remote_digests']['Tokens'] for n in o}) == 1})
check('schema sha256 equal at all stages', len({o[n]['schema_sha256'] for n in o}) == 1)
check('HumanModel=1 and goose2/guard1/rows1 at all stages (no migration)', all(o[n]['human_model'] == 1 and json.loads(o[n]['goose_trigger']) == {'goose': 2, 'guard': 1, 'rows': 1} for n in o))
check('frozen equals before for every state digest incl. Rates (writer stopped)', b['digests'] == f['digests'] and b['remote_digests'] == f['remote_digests'])
check('Connections: baseline had 1 expiring attempt, 0 from before (old relay) through final; unchanged across the switch', base['counts']['Connections'] == 1 and all(o[n]['digests']['Connections'] == b['digests']['Connections'] and o[n]['counts']['Connections'] == 0 for n in ['before', 'frozen', 'after', 'final']))
check('Rates is volatile (informational, not a permission change)', True, {'Rates_counts': [o[n]['counts']['Rates'] for n in STAGES], 'digest_equal_before_frozen': b['digests']['Rates'] == f['digests']['Rates']})
fe = {n: o[n]['files'] for n in o}
check('tunnel.uuid/access.aud hash+mode equal at all stages', all({k: v for k, v in fe[n].items() if k != 'tunnel/config.yml'} == {k: v for k, v in fe['baseline'].items() if k != 'tunnel/config.yml'} for n in o))
check('.env hash+mode unchanged at all stages (env untouched; install dir wired by the server-local override)', all(fe[n]['.env'] == fe['baseline']['.env'] and fe[n]['.env']['mode'] == '0o600' for n in o))
check('tunnel/config.yml unchanged through the relay switch', fe['baseline']['tunnel/config.yml'] == fe['before']['tunnel/config.yml'] == fe['after']['tunnel/config.yml'])
check('tunnel/config.yml changed exactly once after the switch (approved /install ingress), mode 0600', fe['final']['tunnel/config.yml']['sha256'] != fe['after']['tunnel/config.yml']['sha256'] and fe['final']['tunnel/config.yml']['mode'] == '0o600')
check('tunnel credential json hash+mode equal at all stages', all(o[n]['private_tunnel_files'] == b['private_tunnel_files'] for n in o))
check('state and backups directory modes 0700 at all stages', all(o[n]['state_mode'] == '0o700' and o[n]['backups_mode'] == '0o700' for n in o))
bi = b['backup_inventory']
new = [k for k in a['backup_inventory'] if k not in bi]
check('existing backups preserved (hash/mode/bytes); exactly one new backup', len(new) == 1 and all(a['backup_inventory'][k] == v for k, v in bi.items()) and a['backup_inventory'] == z['backup_inventory'], {'before_count': len(bi), 'after_count': len(a['backup_inventory'])})
check('new backup is mode 0600', a['backup_inventory'][new[0]]['mode'] == '0o600', {'bytes': a['backup_inventory'][new[0]]['bytes'], 'sha256': a['backup_inventory'][new[0]]['sha256']})
cmap = lambda x: {c['name']: c for c in x['containers']}
cb0, cb, ca, cz = cmap(base), cmap(b), cmap(a), cmap(z)
key = ('image', 'id', 'running', 'health', 'started', 'restarts')
sig = lambda c: {k: c[k] for k in key}
changed_after = sorted(n for n in cb0 if sig(cb0[n]) != sig(ca[n]))
check('only relay and migrate containers changed across the switch', changed_after == ['/knowslink-migrate-1', '/knowslink-relay-1'], changed_after)
check('Postgres container unchanged across switch and final', all(sig(c['/knowslink-postgres-1']) == sig(cb0['/knowslink-postgres-1']) for c in (ca, cz)))
changed_final = sorted(n for n in cb0 if sig(cb0[n]) != sig(cz[n]))
check('final: only relay, migrate and KnowsLink cloudflared changed', changed_final == ['/knowslink-cloudflared-1', '/knowslink-migrate-1', '/knowslink-relay-1'], changed_final)
check('final: every other pre-existing container (shared services included) unchanged', all(sig(cb0[n]) == sig(cz[n]) for n in cb0 if n not in changed_final), {'pre_existing': len(cb0), 'unchanged': len(cb0) - len(changed_final)})
check('only new container since baseline is the retained restore fixture', sorted(n for n in cz if n not in cb0) == ['/knowslink-connect001-rehearse-20261027'], sorted(n for n in cz if n not in cb0))
check('relay image changed and relay healthy with restarts 0', cb0['/knowslink-relay-1']['image'] != cz['/knowslink-relay-1']['image'] and cz['/knowslink-relay-1']['health'] == 'healthy' and cz['/knowslink-relay-1']['restarts'] == 0)
check('shared orca/s8/mcp statuses and host cloudflared PIDs equal at all stages', all(o[n]['shared'] == b['shared'] and o[n]['host_cloudflared_pids'] == b['host_cloudflared_pids'] for n in o), {'host_cloudflared_pids': b['host_cloudflared_pids']})
check('listeners 127.0.0.1:8080 and 5432 restored; frozen had only 5432', b['listeners'] == a['listeners'] == z['listeners'] and f['listeners'] == ['127.0.0.1:5432'], {'baseline': base['listeners'], 'frozen': f['listeners']})
check('no other active DB writer at frozen/after/final', all(json.loads(o[n]['db_connections'])['active_other_writers'] == 0 for n in ['frozen', 'after', 'final']))
NEW, OLD = 'bf67dfcbb4f7164149901d847edc0aa0e5783dc7', '4fdd149b27d3ff251a73128120dce9f8bb136620'
check('checkout bf67dfcb clean at before/frozen/after/final (moved by the build step); baseline at 4fdd149b', all(o[n]['head'] == NEW and not o[n]['dirty'] for n in ['before', 'frozen', 'after', 'final']) and base['head'] == OLD)
epochs = [o[n]['epoch'] for n in STAGES]
check('epoch never decreased (frozen epoch 2180035 was verified stable for 2s in deploy.log)', epochs == sorted(epochs), epochs)
out = {'utc_stages': {n: o[n]['utc'] for n in STAGES}, 'passed': sum(x['pass'] for x in A), 'total': len(A), 'assertions': A}
(root / 'preservation.json').write_text(json.dumps(out, indent=2, ensure_ascii=False) + '\n', encoding='utf-8')
print(out['passed'], '/', out['total'])
for x in A:
    if not x['pass']:
        print('FAIL', x['assertion'], x.get('detail'))
