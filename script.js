(() => {
  document.addEventListener('DOMContentLoaded', () => init());

  function init() {
    initCopyButtons();
    initTableOfContent();
  }

  function initCopyButtons() {
    let blocks = document.querySelectorAll('.content .highlight');
    blocks.forEach((block) => {
      let btn = document.createElement('div');
      btn.classList.add('btn-copy');
      btn.innerText = 'Copy';
      block.prepend(btn);
      let timeoutId = null;

      btn.addEventListener('click', async () => {
        let content = block.querySelector('pre').innerText;

        try {
          await navigator.clipboard.writeText(content);
          btn.innerText = 'Copied!';
        } catch (err) {
          console.error('Failed to copy:', err);
          btn.innerText = `${err}`;
        }

        if (timeoutId) clearTimeout(timeoutId);
        timeoutId = setTimeout(() => {
          btn.innerText = 'Copy';
        }, 5000);
      });
    });
  }

  function initTableOfContent() {
    const toc$ = document.querySelector('#header-toc');
    const hx$ = Array.from(document.querySelectorAll('.content h2, .content h3'));
    if (!toc$ || hx$.length === 0) return;

    // add data-id to toc li
    const tocLi$ = Array.from(toc$.querySelectorAll('li'));
    tocLi$.forEach((li$) => {
      let href = li$.querySelector('a[href]').href;
      let id = href.split('#')[1];
      li$.dataset.id = id;

      li$.addEventListener('click', (ev) => {
        li$.querySelector('a[href]')?.click();
        ev.stopPropagation();
      });
    });

    let lastTimeoutId;
    let lastProcessTime = 0;

    const observer = new IntersectionObserver(
      (entries) => {
        let hasActive = false;
        entries.forEach((entry) => {
          if (entry.isIntersecting) hasActive = true;
        });
        if (!hasActive) return;

        if (Date.now() - lastProcessTime > 500) {
          processObserver(entries);
          lastProcessTime = Date.now();
          return;
        }

        if (lastTimeoutId) {
          clearTimeout(lastTimeoutId);
          lastTimeoutId = undefined;
        }
        lastTimeoutId = setTimeout(() => {
          processObserver(entries);
          lastProcessTime = Date.now();
        }, 500);
      },
      {threshold: 0.5},
    );
    hx$.forEach((hx) => {
      observer.observe(hx);
    });

    let lastActiveId;

    function processObserver(entries) {
      // check if any h2 or h3 is intersecting
      let [h2, h3] = [false, false];
      let mapActive = {};

      // set h2, h3
      entries.forEach((entry) => {
        let isActive = entry.isIntersecting;
        if (isActive && entry.target.tagName === 'H2') h2 = true;
        if (isActive && entry.target.tagName === 'H3') h3 = true;
      });

      // set mapActive
      entries.forEach((entry) => {
        let id = entry.target.id;
        let isActive = entry.isIntersecting;
        if (entry.target.tagName === 'H2' && h3) {
          isActive = false;
        }
        if (isActive) {
          lastActiveId = id;
          mapActive[id] = true;
        }
      });

      // if there is no active li, use lastActiveId
      if (Object.keys(mapActive).length === 0) {
        if (lastActiveId) mapActive[lastActiveId] = true;
      }

      // update toc li.active
      tocLi$.forEach((li) => {
        let id = li.dataset.id;
        if (mapActive[id]) {
          if (!li.classList.contains('active')) li.classList.add('active');
        } else {
          if (li.classList.contains('active')) li.classList.remove('active');
        }
      });
    }
  }
})();
