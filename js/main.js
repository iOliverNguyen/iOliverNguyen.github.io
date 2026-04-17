/**
 * Sets up Justified Gallery.
 */
if (!!$.prototype.justifiedGallery) {
  let options = {
    rowHeight: 140,
    margins: 4,
    lastRow: 'justify',
  };
  $('.article-gallery').justifiedGallery(options);
}

$(document).ready(function () {

  /**
   * Shows the responsive navigation menu on mobile.
   */
  $('#header > #nav > ul > .icon').click(function () {
    $('#header > #nav > ul').toggleClass('responsive');
  });


  /**
   * Controls the different versions of  the menu in blog post articles
   * for Desktop, tablet and mobile.
   */
  if ($('.post').length) {
    let menu = $('#menu');
    let nav = $('#menu > #nav');
    let menuIcon = $('#menu-icon, #menu-icon-tablet');

    /**
     * Display the menu on hi-res laptops and desktops.
     */
    if ($(document).width() >= 1440) {
      menu.css('visibility', 'visible');
      menuIcon.addClass('active');
    }

    /**
     * Display the menu if the menu icon is clicked.
     */
    // menuIcon.click(function() {
    //   if (menu.css('visibility') === 'hidden') {
    //     menu.css('visibility', 'visible');
    //     menuIcon.addClass('active');
    //   } else {
    //     menu.css('visibility', 'hidden');
    //     menuIcon.removeClass('active');
    //   }
    //   return false;
    // });

    /**
     * Add a scroll listener to the menu to hide/show the navigation links.
     */
    {
      let lastScrollTop = 0;
      $(window).on('scroll', function () {
        let topDistance = $(window).scrollTop();
        if (topDistance > lastScrollTop + 10) {
          // downscroll -> hide menu
          $('#header-post #menu-icon').css('opacity', '0');
          $('#header-post #menu-icon-tablet').css('opacity', '0');
        } else if (topDistance <= lastScrollTop) {
          // upscroll -> show menu
          $('#header-post #menu-icon').css('opacity', '1');
          $('#header-post #menu-icon-tablet').css('opacity', '1');
        }
        lastScrollTop = topDistance;
      });
    }

    // if (menu.length) {
    //   $(window).on('scroll', function() {
    //     let topDistance = menu.offset().top;
    //
    //     // hide only the navigation links on desktop
    //     if (!nav.is(':visible') && topDistance < 50) {
    //       nav.show();
    //     } else if (nav.is(':visible') && topDistance > 100) {
    //       nav.hide();
    //     }
    //
    //     // on tablet, hide the navigation icon as well and show a 'scroll to top
    //     // icon' instead
    //     if ( ! $( '#menu-icon' ).is(':visible') && topDistance < 50 ) {
    //       $('#menu-icon-tablet').show();
    //       $('#top-icon-tablet').hide();
    //     } else if (! $( '#menu-icon' ).is(':visible') && topDistance > 100) {
    //       $('#menu-icon-tablet').hide();
    //       $('#top-icon-tablet').show();
    //     }
    //   });
    // }

    /**
     * Show mobile navigation menu after scrolling upwards,
     * hide it again after scrolling downwards.
     */

    let justClickedTOC = false;
    if ($('#footer-post').length) {
      $(window).on('resize', () => {
        // resize -> show menu
        $('#footer-post').css('display', '');
      });

      let lastScrollTop = 0;
      $(window).on('scroll', function () {
        let topDistance = $(window).scrollTop();
        if (topDistance > lastScrollTop + 10) {
          // downscroll -> hide menu, but only if it's not from clicking TOC
          if (!justClickedTOC) $('#footer-post').css('display', 'none');

        } else if (topDistance <= lastScrollTop) {
          // upscroll -> show menu
          $('#footer-post').css('display', '');
        }
        lastScrollTop = topDistance;
        justClickedTOC = false;

        // close all submenu's on scroll
        $('#nav-footer').hide();
        $('#toc-footer').hide();
        $('#share-footer').hide();

        // show a 'navigation' icon when close to the top of the page,
        // otherwise show a 'scroll to the top' icon
        // if (topDistance < 50) {
        //   $('#actions-footer > #top').hide();
        // } else if (topDistance > 100) {
        //   $('#actions-footer > #top').show();
        // }
      });
    }

    if ($('#toc-footer #TableOfContents').length) {
      // if [class=alt] is present, hide previous content
      $('#toc-footer #TableOfContents .alt').map((_, elem) => {
        let list = [];
        let prev = elem.previousSibling;
        while (prev) {
          list.push(prev);
          prev = prev.previousSibling;
        }
        for (let x of list) { x.remove() }
      });

      // always show footer-post after navigating
      $('#toc-footer #TableOfContents a').click(() => {
        justClickedTOC = true;
      });
    }
  }
});
