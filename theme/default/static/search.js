(function() {
  var index = null;
  var container = document.getElementById('search-container');
  var input = document.getElementById('search-input');
  var results = document.getElementById('search-results');
  if (!container || !input) return;

  document.addEventListener('keydown', function(e) {
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      e.preventDefault();
      if (container.style.display === 'none') {
        container.style.display = 'block';
        input.focus();
      } else {
        container.style.display = 'none';
      }
    }
    if (e.key === 'Escape') { container.style.display = 'none'; }
  });

  input.addEventListener('focus', function() {
    if (index) return;
    var baseURL = (document.querySelector('meta[name="base-url"]') || {}).content || '';
    fetch(baseURL + '/search-index.json')
      .then(function(r) { return r.json(); })
      .then(function(data) { index = data; });
  });

  function search(query) {
    if (!index || !query) return [];
    var q = query.toLowerCase();
    var scored = [];
    for (var i = 0; i < index.length; i++) {
      var entry = index[i];
      var score = 0;
      var fields = [entry.title, entry.titleKo, entry.brief].concat(entry.tags || []).concat(entry.aliases || []);
      for (var j = 0; j < fields.length; j++) {
        if (fields[j] && fields[j].toLowerCase().indexOf(q) !== -1) {
          score += (j === 0 || j === 1) ? 10 : 1;
        }
      }
      if (score > 0) scored.push({ entry: entry, score: score });
    }
    scored.sort(function(a, b) { return b.score - a.score; });
    return scored.slice(0, 10);
  }

  input.addEventListener('input', function() {
    var matches = search(input.value);
    results.innerHTML = '';
    var baseURL = (document.querySelector('meta[name="base-url"]') || {}).content || '';
    for (var i = 0; i < matches.length; i++) {
      var m = matches[i].entry;
      var li = document.createElement('li');
      var a = document.createElement('a');
      a.href = baseURL + '/pages/' + encodeURIComponent(m.name);
      a.textContent = m.titleKo || m.title;
      if (m.brief) {
        var span = document.createElement('span');
        span.textContent = ' — ' + m.brief;
        span.style.cssText = 'color:var(--c-text-muted);font-size:var(--text-s)';
        a.appendChild(span);
      }
      li.appendChild(a);
      results.appendChild(li);
    }
  });
})();
