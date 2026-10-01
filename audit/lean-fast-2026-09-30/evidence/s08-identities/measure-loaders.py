from pathlib import Path
import subprocess,os,re,json,statistics
p=Path('/tmp/mpress-s08-identities-20261001');root=Path('/var/home/lea/projects/mpress-audit-20260930');rows=[]
for corpus,site in [('legacy',p/'legacy-site'),('corrected',Path('/tmp/mpress-audit-20260930/wails-corpus/docs/mpress/site'))]:
 for iteration in range(3):
  for kind in ['before','after']:
   r=subprocess.run([str(p/(kind+'-knowledge.test')),'-test.run=^$','-test.bench=^BenchmarkExplicitKnowledgeLoadAll$','-test.benchmem','-test.benchtime=1x','-test.count=1'],env=dict(os.environ,MPRESS_KNOWLEDGE_CORPUS=str(site)),cwd=root,capture_output=True,text=True)
   (p/f'bench-{corpus}-{kind}-{iteration}.log').write_text(r.stdout+r.stderr);assert r.returncode==0,r.stdout+r.stderr
   m=re.search(r'BenchmarkExplicitKnowledgeLoadAll-\d+\s+1\s+(\d+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op',r.stdout);assert m,r.stdout
   rows.append({'corpus':corpus,'candidate':kind,'iteration':iteration,'ns':int(m[1]),'allocatedBytes':int(m[2]),'allocations':int(m[3])})
summary={}
for corpus in ['legacy','corrected']:
 summary[corpus]={}
 for kind in ['before','after']:
  subset=[r for r in rows if r['corpus']==corpus and r['candidate']==kind]
  summary[corpus][kind]={k:statistics.median(r[k] for r in subset) for k in ['ns','allocatedBytes','allocations']}
(p/'loader-cost.json').write_text(json.dumps({'baseline':'c4b8ddf','alternatingRunsPerCandidatePerCorpus':3,'sameInputPerComparison':True,'cpuProfilingDuringTiming':False,'memoryMetric':'Total Go allocation bytes/op, not peak RSS','samples':rows,'medians':summary},indent=2)+'\n');print(json.dumps(summary),flush=True)
