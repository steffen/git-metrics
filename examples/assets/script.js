// script.js
// Load the metrics output, segment it into sections, and provide synchronized navigation.

// Supported outputs. First entry loads initially.
const OUTPUT_FILES = ['outputs/git.txt', 'outputs/linux.txt', 'outputs/chromium.txt'];
let currentOutputIndex = 0;
function currentOutputFile(){ return OUTPUT_FILES[currentOutputIndex]; }
// Added: alignment factor for vertical positioning of focused section (0 = top, 0.5 = center)
const FOCUS_VERTICAL_ALIGN = 0.35; // 35% from top gives a balanced look

// SECTION DEFINITIONS ARE NOW EXTERNALIZED IN sections.md
// Parsed at runtime so that documentation and UI stay in sync.
let SECTION_DEFINITIONS = [];

async function loadSectionDefinitions(){
  const res = await fetch('sections.md'); // relative to assets/ directory (index.html in examples/)
  if(!res.ok){
    throw new Error('Failed to load sections.md: '+res.status);
  }
  const md = await res.text();
  SECTION_DEFINITIONS = parseSectionsMarkdown(md);
}

function parseSectionsMarkdown(md){
  const lines = md.split(/\r?\n/);
  const sections = [];
  let current = null;
  let descriptionLines = [];

  const commit = () => {
    if(current){
      const description = descriptionLines.join(' ').trim();
      sections.push(buildSectionDefinition(current, description));
    }
    current = null;
    descriptionLines = [];
  };

  for(const rawLine of lines){
    const line = rawLine.trimEnd();
    if(/^# /.test(line)){ // New top-level section
      commit();
      current = { header: line.replace(/^#\s+/,'').trim(), title: null };
      continue;
    }
    if(/^## /.test(line)){ // Title line
      if(!current){
        // If a secondary heading appears before a primary, start an implicit section
        current = { header: line.replace(/^##\s+/,'').trim(), title: null };
      }
      current.title = line.replace(/^##\s+/,'').trim();
      continue;
    }
    if(line.startsWith('#')){ // Any other heading depth ends previous section
      commit();
      continue;
    }
    if(line.trim().length === 0){
      // blank line just separates paragraphs; retain paragraph spacing subtly by adding space later
      continue;
    }
    descriptionLines.push(line.trim());
  }
  commit();
  return sections;
}

function buildSectionDefinition(parsed, description){
  const header = parsed.header || parsed.title || '';
  const title = parsed.title || header; // fallback
  const id = deriveSectionId(header);
  const match = deriveMatchRegex(id, header);
  return {
    id,
    title,
    match,
    explain: () => description
  };
}

function deriveSectionId(header){
  const norm = header.toLowerCase();
  const map = {
    'run': 'run',
    'run and repository': 'run-and-repository',
    'repository': 'repository',
    'historic and estimated growth': 'growth',
    'historic & estimated growth': 'growth',
    'rate of changes': 'rate',
    'largest directories': 'largest-dirs',
    'largest files': 'largest-files',
    'largest file extensions': 'largest-ext',
    'authors with most commits': 'top-authors',
    'committers with most commits': 'top-committers',
    'footer': 'footer'
  };
  return map[norm] || norm.replace(/[^a-z0-9]+/g,'-').replace(/^-|-$/g,'');
}

function deriveMatchRegex(id, header){
  if(id === 'footer') return /^Finished in /i; // footer remains matched by runtime summary line
  if(id === 'run-and-repository') return /^RUN /i; // anchor combined section to first RUN header in output
  // Escape regex special chars in header text and expect a trailing space in output line (as current format)
  const escaped = header.replace(/[.*+?^${}()|[\]\\]/g,'\\$&');
  return new RegExp('^' + escaped + ' ', 'i');
}

// Mapping of section identifiers to one or more Git commands (or derivations) used to produce the data.
// Placeholders:
//   <default_branch>  – resolved default branch (e.g. main / master)
//   <year>            – iterated year during growth calculation loop
//   <file>            – each candidate file when determining largest files
const GIT_COMMANDS_BY_SECTION = {
  'run': [
    'git version'
  ],
  'repository': [
    'git remote get-url origin',
    'git rev-parse --short HEAD',
    'git show -s --format=%cD <hash>',
    'git rev-list --max-parents=0 HEAD --format=%cD'
  ],
  // Combined (Alternative 1) section covering both run + repository metadata
  'run-and-repository': [
    'git version',
    'git remote get-url origin',
    'git rev-parse --short HEAD',
    'git show -s --format=%cD <hash>',
    'git rev-list --max-parents=0 HEAD --format=%cD'
  ],
  'growth': [
    // For every <year> in the range first_commit_year..current_year
    `git rev-list --objects --all --before <year>-01-01 --after <year>-12-31 | git cat-file --batch-check='%(objecttype) %(objectname) %(objectsize:disk) %(rest)'`
  ],
  'top-authors': [
    'git log --all --format=%an|%cn|%cd --date=format:%Y'
  ],
  'top-committers': [
    'git log --all --format=%an|%cn|%cd --date=format:%Y'
  ],
  'rate': [
    'git remote show origin   # read default branch (HEAD branch line)',
    'git log <default_branch> --format=%ct|%P --reverse'
  ],
  'largest-dirs': [
    'git ls-tree -r --name-only <default_branch>'
  ],
  'largest-files': [
    'git log -1 --format=%cD -- <file>   # per listed file'
  ],
  'largest-ext': [
    // Extension statistics are derived from the blob listing produced by growth commands
    '(derived from growth blob inventory)'
  ],
  'footer': []
};

const outputPane = document.getElementById('outputPane');
const explanationPane = document.getElementById('explanationPane');
const repoBadgeEl = document.getElementById('repoBadge'); // deprecated visual usage (hidden)
const outputContentEl = document.getElementById('outputContent');
let outputTitleEl = null; // large title on right side

let sections = []; // [{id, startLine, endLine, explanationEl}]
let activeIndex = 0;
let outputIndicatorEl = null;

function deriveRepoName(filePath){
  const base = filePath.split('/').pop() || filePath;
  return base.replace(/\.txt$/,'');
}

async function loadOutputFile(preserveSectionId){
  const file = currentOutputFile();
  const res = await fetch(file);
  const text = await res.text();
  const lines = text.split(/\n/);

  // Clear previous content
  outputContentEl.innerHTML = '';
  explanationPane.innerHTML = '';
  sections = [];

  // Update (hidden) repository badge for accessibility only
  if(repoBadgeEl){
    repoBadgeEl.textContent = deriveRepoName(file);
  }

  // Reset definitions line indices before reuse
  SECTION_DEFINITIONS.forEach(d=>{ delete d.lineIndex; delete d.endLine; });

  // Build section boundaries
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

  // Render output lines
  const frag = document.createDocumentFragment();
  lines.forEach((ln,i)=>{
    const div = document.createElement('div');
    div.textContent = ln || '\u200b';
    div.className = 'masked-line';
    div.dataset.line = i;
    frag.appendChild(div);
  });
  outputContentEl.appendChild(frag);

  // Explanation sections
  const expFrag = document.createDocumentFragment();
  sortedDefs.forEach(def => {
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

  // Attach Git command displays after section elements exist
  attachGitCommands();

  buildAnchorBar();

  // Ensure output title exists (large heading on right side)
  if(!outputTitleEl){
    outputTitleEl = document.createElement('h1');
    outputTitleEl.className = 'output-title';
  }
  outputTitleEl.textContent = deriveRepoName(file);
  // Prepend so it appears above anchor bar
  explanationPane.prepend(outputTitleEl);

  // Determine active index to restore
  let restoreIndex = 0;
  if(preserveSectionId){
    const idx = sections.findIndex(s=>s.id===preserveSectionId);
    if(idx !== -1) restoreIndex = idx;
  }
  updateActiveSection(restoreIndex, false);
}

// Add a bottom-left list of Git commands used for each section (if any)
function attachGitCommands(){
  sections.forEach(sectionEntry => {
    const commands = GIT_COMMANDS_BY_SECTION[sectionEntry.id];
    if(!commands || commands.length === 0){
      return; // Nothing to show
    }
    const container = document.createElement('div');
    container.className = 'git-command-display';
    container.setAttribute('aria-label','Git commands used to generate this section');
    // Optional label (kept visually subtle via CSS)
    const label = document.createElement('div');
    label.className = 'git-command-label';
    label.textContent = 'Git commands';
    container.appendChild(label);
    // Add wrap opportunities after pipe and logical AND for long commands
    const insertWrapHints = (text) => text
      .replace(/\|/g, '|\u200b')
      .replace(/&&/g, '&&\u200b');
    commands.forEach(commandString => {
      const codeElement = document.createElement('code');
      codeElement.textContent = insertWrapHints(commandString);
      container.appendChild(codeElement);
    });
    sectionEntry.explanationEl.appendChild(container);
  });
}

function buildAnchorBar(){
  const bar = document.createElement('div');
  bar.className='anchor-links';
  // Output indicator
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
  // Section navigation (Up/Down)
  if(e.key==='ArrowDown'){ e.preventDefault(); updateActiveSection(Math.min(activeIndex+1, sections.length-1)); }
  if(e.key==='ArrowUp'){ e.preventDefault(); updateActiveSection(Math.max(activeIndex-1, 0)); }
  // Output switching (Right/Left) preserving section id
  if(e.key==='ArrowRight'){
    if(currentOutputIndex < OUTPUT_FILES.length - 1){
      e.preventDefault();
      const preserveId = sections[activeIndex]?.id;
      currentOutputIndex++;
      loadOutputFile(preserveId).catch(err=>{ outputPane.textContent = 'Failed to load output: '+err; });
    }
  }
  if(e.key==='ArrowLeft'){
    if(currentOutputIndex > 0){
      e.preventDefault();
      const preserveId = sections[activeIndex]?.id;
      currentOutputIndex--;
      loadOutputFile(preserveId).catch(err=>{ outputPane.textContent = 'Failed to load output: '+err; });
    }
  }
}

function refreshAnchors(){
  const links = document.querySelectorAll('.anchor-links a');
  links.forEach((a,i)=>a.classList.toggle('active', i===activeIndex));
}

// Initialization sequence: load section definitions first, then load the current output.
init();
async function init(){
  try {
    await loadSectionDefinitions();
    await loadOutputFile();
    window.addEventListener('keydown', onKey);
    outputPane.addEventListener('scroll', onManualScroll);
    initNavHintBehavior();
  } catch(err){
    outputPane.textContent = 'Failed to initialize: ' + err;
    console.error(err);
  }
}

// --- Navigation Hint Show/Hide Logic ---------------------------------------------------------
const NAV_HINT_INITIAL_VISIBLE_MS = 4000; // how long to keep it visible on first load
const NAV_HINT_INACTIVITY_MS = 2500; // hide after this long without mouse movement
let navHintEl = document.querySelector('.nav-hint');
let navHintHideTimeout = null;
let navHintRecentlyShown = false;

function initNavHintBehavior(){
  if(!navHintEl) return;
  // Ensure visible at start
  showNavHint(false);
  // Schedule initial hide
  navHintHideTimeout = setTimeout(()=>hideNavHint(), NAV_HINT_INITIAL_VISIBLE_MS);
  // Listen for mouse movement to re-show
  window.addEventListener('mousemove', onUserActivityForNavHint, { passive:true });
  window.addEventListener('keydown', onUserActivityForNavHint, { passive:true }); // keyboard navigation also reveals
}

function onUserActivityForNavHint(){
  if(!navHintEl) return;
  // If already scheduled, reset timer
  if(navHintHideTimeout){
    clearTimeout(navHintHideTimeout);
  }
  // Only re-show if currently hidden OR not shown very recently (debounce flicker)
  if(navHintEl.classList.contains('is-hidden') || !navHintRecentlyShown){
    showNavHint(true);
  }
  navHintHideTimeout = setTimeout(()=>hideNavHint(), NAV_HINT_INACTIVITY_MS);
}

function showNavHint(animated){
  if(!navHintEl) return;
  navHintEl.classList.remove('is-hidden');
  navHintRecentlyShown = true;
  // Cooldown to prevent rapid toggling animations
  setTimeout(()=>{ navHintRecentlyShown = false; }, 800);
}

function hideNavHint(){
  if(!navHintEl) return;
  navHintEl.classList.add('is-hidden');
}
