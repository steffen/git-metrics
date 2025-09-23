// script.js
// Load the metrics output, segment it into sections, and provide synchronized navigation.

const OUTPUT_FILE = 'outputs/git.txt';
// Added: alignment factor for vertical positioning of focused section (0 = top, 0.5 = center)
const FOCUS_VERTICAL_ALIGN = 0.35; // 35% from top gives a balanced look

const SECTION_DEFINITIONS = [
  { id: 'run', title: 'Run Metadata', match: /^RUN /i, explain: () => `General metadata about when and where the report was generated (start time, host machine, versions). Useful to contextualize measurements and reproduce runs.`},
  { id: 'repository', title: 'Repository Info', match: /^REPOSITORY /i, explain: () => `Origin information: repository path, remote URL, most recent commit and age. Helps anchor the report to a specific revision set.`},
  { id: 'growth', title: 'Historic & Estimated Growth', match: /^HISTORIC & ESTIMATED GROWTH /i, explain: () => `Shows yearly totals of core Git objects (commits, trees, blobs) along with on-disk size. Past years are actual; rows with ^ are current totals; rows with * are projections extrapolated from recent growth.`},
  { id: 'rate', title: 'Rate of Changes', match: /^RATE OF CHANGES /i, explain: () => `Focuses on commit cadence to the default branch. P95/P99/P100 peaks per day/hour/minute reveal burstiness and scaling of integration workflow.`},
  { id: 'largest-dirs', title: 'Largest Directories', match: /^LARGEST DIRECTORIES /i, explain: () => `Identifies directories contributing ≥1% of repository storage. Highlights translation files, tests, docs, and core source areas for optimization or pruning.`},
  { id: 'largest-files', title: 'Largest Files', match: /^LARGEST FILES /i, explain: () => `Top individual files by cumulative blob storage, signaling hotspots for size bloat and potential candidates for history rewriting or splitting.`},
  { id: 'largest-ext', title: 'Largest File Extensions', match: /^LARGEST FILE EXTENSIONS /i, explain: () => `Distribution of blob count and size by file extension. Useful to see language / artifact composition and track shifts over time.`},
  { id: 'top-authors', title: 'Authors With Most Commits', match: /^AUTHORS WITH MOST COMMITS /i, explain: () => `Per-year top authors by authored commits plus totals. Shows contributor concentration and evolution of community participation.`},
  { id: 'top-committers', title: 'Committers With Most Commits', match: /^COMMITTERS WITH MOST COMMITS /i, explain: () => `Committer stats (who integrated patches). High centralization can indicate a gatekeeping pattern or strong maintainer oversight.`},
  { id: 'footer', title: 'Footer / Summary', match: /^Finished in /i, explain: () => `Runtime performance of the metrics tool itself (execution time, memory footprint).`}
];

const outputPane = document.getElementById('outputPane');
const explanationPane = document.getElementById('explanationPane');

let sections = []; // [{id, startLine, endLine, element, explanationEl}]
let activeIndex = 0;

async function loadOutput() {
  const res = await fetch(OUTPUT_FILE);
  const text = await res.text();
  const lines = text.split(/\n/);

  // Build an array of section boundaries.
  SECTION_DEFINITIONS.forEach(def => {
    const idx = lines.findIndex(l => def.match.test(l));
    if (idx !== -1) def.lineIndex = idx; else def.lineIndex = Infinity;
  });
  const sortedDefs = SECTION_DEFINITIONS.filter(d => d.lineIndex !== Infinity).sort((a,b)=>a.lineIndex-b.lineIndex);
  for (let i=0;i<sortedDefs.length;i++) {
    const current = sortedDefs[i];
    const next = sortedDefs[i+1];
    current.endLine = (next ? next.lineIndex : lines.length) - 1;
  }

  // Render output lines with data-line attributes for masking.
  const frag = document.createDocumentFragment();
  lines.forEach((ln, i) => {
    const div = document.createElement('div');
    div.textContent = ln || '\u200b';
    div.className = 'masked-line';
    div.dataset.line = i;
    frag.appendChild(div);
  });
  outputPane.appendChild(frag);

  // Create explanation sections.
  const expFrag = document.createDocumentFragment();
  sortedDefs.forEach((def, i) => {
    const sec = { id: def.id, def, startLine: def.lineIndex, endLine: def.endLine };
    const wrap = document.createElement('section');
    wrap.className = 'section-explanation';
    wrap.id = 'exp-' + def.id;
    wrap.innerHTML = `\n      <span class="section-anchor" id="anchor-${def.id}"></span>\n      <h2>${def.title}</h2>\n      <p>${def.explain()}</p>\n      <div class="details" data-section-details></div>\n    `;
    sec.explanationEl = wrap;
    sections.push(sec);
    expFrag.appendChild(wrap);
  });
  explanationPane.appendChild(expFrag);

  buildAnchorBar();
  updateActiveSection(0, false);
  window.addEventListener('keydown', onKey);
  outputPane.addEventListener('scroll', onManualScroll);
}

function buildAnchorBar(){
  const bar = document.createElement('div');
  bar.className='anchor-links';
  sections.forEach((s,i)=>{
    const a=document.createElement('a');
    a.href='#anchor-'+s.id; a.textContent=(i+1)+'. '+s.def.title.split(' ')[0];
    a.addEventListener('click',e=>{e.preventDefault();updateActiveSection(i,true);});
    bar.appendChild(a);
  });
  explanationPane.prepend(bar);
}

let syncing = false; // prevents feedback loops
function onManualScroll(){
  if(syncing) return; // ignore programmatic scroll
  const scrollTop = outputPane.scrollTop;
  const lineHeight = outputPane.querySelector('.masked-line')?.offsetHeight || 14;
  const firstVisibleLine = Math.floor(scrollTop / lineHeight);
  const idx = sections.findIndex(s => firstVisibleLine >= s.startLine && firstVisibleLine <= s.endLine);
  if (idx !== -1 && idx !== activeIndex){
    updateActiveSection(idx,false);
  }
}

function updateActiveSection(newIndex, smooth=true){
  if(newIndex <0 || newIndex >= sections.length) return;
  activeIndex = newIndex;
  const sec = sections[newIndex];

  // Highlight explanation
  sections.forEach(s => s.explanationEl.classList.toggle('active', s===sec));

  // Compute focus range lines
  const focusStart = sec.startLine;
  const focusEnd = sec.endLine;

  // Apply dimming
  outputPane.classList.add('dimmed');
  document.querySelectorAll('.masked-line').forEach(el=>{
    const ln = +el.dataset.line;
    el.classList.toggle('focus-range', ln>=focusStart && ln<=focusEnd);
  });

  // Scroll both panes
  const targetLineEl = outputPane.querySelector(`.masked-line[data-line='${focusStart}']`);
  if(targetLineEl){
    syncing = true;
    // CUSTOM SCROLL LOGIC (replaces scrollIntoView start): center (approx) the focused section
    const containerHeight = outputPane.clientHeight;
    const sectionMidLine = focusStart + (Math.max(0, focusEnd - focusStart) / 2);
    const midLineEl = outputPane.querySelector(`.masked-line[data-line='${Math.round(sectionMidLine)}']`) || targetLineEl;
    const desiredScrollTop = (midLineEl.offsetTop || 0) - containerHeight * FOCUS_VERTICAL_ALIGN;
    const clampedScrollTop = Math.max(0, Math.min(desiredScrollTop, outputPane.scrollHeight - containerHeight));
    outputPane.scrollTo({ top: clampedScrollTop, behavior: smooth? 'smooth':'instant' });
    const y = sec.explanationEl.offsetTop;
    explanationPane.scrollTo({top:y, behavior: smooth?'smooth':'instant'});
    setTimeout(()=>{syncing=false;}, 400);
  }
  refreshAnchors();
}

function onKey(e){
  if(e.key==='ArrowRight') { e.preventDefault(); updateActiveSection(Math.min(activeIndex+1, sections.length-1)); }
  if(e.key==='ArrowLeft') { e.preventDefault(); updateActiveSection(Math.max(activeIndex-1, 0)); }
}

function refreshAnchors(){
  const links = document.querySelectorAll('.anchor-links a');
  links.forEach((a,i)=>a.classList.toggle('active', i===activeIndex));
}

loadOutput().catch(err=>{
  outputPane.textContent = 'Failed to load output: '+err;
});
