document.addEventListener("DOMContentLoaded", function (event) {
    let path = window.location.pathname;
    let navs = document.getElementsByClassName('gw-nav');
    Array.prototype.filter.call(navs, function (nav) {
        if (nav.getAttribute("href") === path) {
            nav.classList.add("active");
        }
    });

    document.getElementById("sidebarToggle").addEventListener("click", function(e){
        e.preventDefault();
        let b = document.getElementsByTagName("BODY")[0];
        if (b.classList.contains("sb-sidenav-toggled")) {
            b.classList.remove("sb-sidenav-toggled");
        } else {
            b.classList.add("sb-sidenav-toggled");
        }
    });
});
