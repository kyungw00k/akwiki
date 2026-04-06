(function() {
  var index = null;
  var input = document.getElementById('search-input');
  var results = document.getElementById('search-results');
  if (!input || !results) return;

  var baseURL = (document.querySelector('meta[name="base-url"]') || {}).content || '';
  var selected = -1;

  function loadIndex() {
    if (index) return Promise.resolve(index);
    return fetch(baseURL + '/search-index.json')
      .then(function(r) { return r.json(); })
      .then(function(data) { index = data; return data; });
  }

  input.addEventListener('focus', loadIndex);

  function escapeHtml(s) {
    return s.replace(/[&<>"']/g, function(c) {
      return { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c];
    });
  }

  function highlight(text, query) {
    if (!text || !query) return escapeHtml(text || '');
    var lower = text.toLowerCase();
    var q = query.toLowerCase();
    var idx = lower.indexOf(q);
    if (idx === -1) return escapeHtml(text);
    return escapeHtml(text.slice(0, idx)) +
           '<mark>' + escapeHtml(text.slice(idx, idx + query.length)) + '</mark>' +
           escapeHtml(text.slice(idx + query.length));
  }

  function search(query) {
    if (!index || !query) return [];
    var q = query.toLowerCase();
    var scored = [];
    for (var i = 0; i < index.length; i++) {
      var entry = index[i];
      var score = 0;
      var matchedAlias = null;
      // Title fields (high weight)
      if (entry.title && entry.title.toLowerCase().indexOf(q) !== -1) score += 10;
      if (entry.titleKo && entry.titleKo.toLowerCase().indexOf(q) !== -1) score += 10;
      // Aliases (medium weight, captured for display)
      if (entry.aliases) {
        for (var k = 0; k < entry.aliases.length; k++) {
          if (entry.aliases[k].toLowerCase().indexOf(q) !== -1) {
            score += 5;
            if (!matchedAlias) matchedAlias = entry.aliases[k];
          }
        }
      }
      // Tags
      if (entry.tags) {
        for (var k2 = 0; k2 < entry.tags.length; k2++) {
          if (entry.tags[k2].toLowerCase().indexOf(q) !== -1) score += 2;
        }
      }
      // Brief (lower weight)
      if (entry.brief && entry.brief.toLowerCase().indexOf(q) !== -1) score += 1;

      if (score > 0) scored.push({ entry: entry, score: score, matchedAlias: matchedAlias });
    }
    scored.sort(function(a, b) { return b.score - a.score; });
    return scored.slice(0, 10);
  }

  function render(matches, query) {
    results.innerHTML = '';
    selected = matches.length > 0 ? 0 : -1;
    for (var i = 0; i < matches.length; i++) {
      var m = matches[i].entry;
      var alias = matches[i].matchedAlias;
      var li = document.createElement('li');
      li.setAttribute('role', 'option');
      if (i === 0) li.classList.add('selected');
      var a = document.createElement('a');
      a.href = baseURL + '/pages/' + encodeURIComponent(m.name);
      var title = m.titleKo || m.title;
      var html = highlight(title, query);
      if (alias && alias !== title) {
        html += '<span class="sub"> (' + highlight(alias, query) + ')</span>';
      }
      a.innerHTML = html;
      li.appendChild(a);
      results.appendChild(li);
    }
  }

  function updateSelection() {
    var items = results.querySelectorAll('li');
    for (var i = 0; i < items.length; i++) {
      items[i].classList.toggle('selected', i === selected);
    }
    if (selected >= 0 && items[selected]) {
      items[selected].scrollIntoView({ block: 'nearest' });
    }
  }

  input.addEventListener('input', function() {
    loadIndex().then(function() {
      render(search(input.value), input.value);
    });
  });

  document.addEventListener('click', function(e) {
    if (!input.contains(e.target) && !results.contains(e.target)) {
      results.innerHTML = '';
    }
  });

  input.addEventListener('keydown', function(e) {
    var items = results.querySelectorAll('li');
    if (e.key === 'Escape') {
      input.value = '';
      results.innerHTML = '';
      input.blur();
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (items.length === 0) return;
      selected = (selected + 1) % items.length;
      updateSelection();
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (items.length === 0) return;
      selected = (selected - 1 + items.length) % items.length;
      updateSelection();
    } else if (e.key === 'Enter') {
      if (selected >= 0 && items[selected]) {
        e.preventDefault();
        var link = items[selected].querySelector('a');
        if (link) window.location.href = link.href;
      }
    }
  });
})();
