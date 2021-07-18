; (function($, window, document, undefined) {
    var pluginName = "jqueryAccordionMenu";
    var defaults = {
        speed: 300,
        showDelay: 0,
        hideDelay: 0,
        singleOpen: true,
        clickEffect: true,
        expandAllFound: false,
        redirectChildFirst: false
    };
    function Plugin(element, options) {
        this.element = element;
        this.settings = $.extend({}, defaults, options);
        this._defaults = defaults;
        this._name = pluginName;
        this.init();
    };
    $.extend(Plugin.prototype, {
        init: function() {
            this.openSubmenu();
            this.submenuIndicators();
            if (this.settings.clickEffect) {
                this.addClickEffect();
            }
            this.filterList();
        },
        openSubmenu: function() {
            var settings = this.settings;
            $(this.element).children("ul").find("li").bind("click touchstart", function(e) {
                e.stopPropagation();
                e.preventDefault();
                var submenu = $(this).children(".submenu");
                if (submenu.length > 0) {
                    if (submenu.css("display") == "none") {
                        submenu.delay(settings.showDelay).slideDown(settings.speed);
                        submenu.siblings("a").addClass("submenu-indicator-minus");
                        if (settings.singleOpen) {
                            var otherSubmenu = $(this).siblings().children(".submenu");
                            otherSubmenu.slideUp(settings.speed);
                            otherSubmenu.siblings("a").removeClass("submenu-indicator-minus")
                        }
                        return false
                    } 
                    submenu.delay(settings.hideDelay).slideUp(settings.speed);
                    var otherA = submenu.siblings("a");
                    if (otherA.hasClass("submenu-indicator-minus")) {
                        otherA.removeClass("submenu-indicator-minus");
                    }
                }
                if(settings.redirectChildFirst){
                    window.location.href = $(this).children("a").attr("href")
                }
            })
        },
        submenuIndicators: function() {
            var submenu = $(this.element).find(".submenu");
            if (submenu.length > 0) {
                submenu.siblings("a").append("<span class='submenu-indicator'>+</span>")
            }
        },
        addClickEffect: function() {
            var that = this;
            $(this.element).find("a").bind("click touchstart", function(e) {
                $(that.element).find(".ink").remove();
                if ($(this).children(".ink").length === 0) {
                    $(this).prepend("<span class='ink'></span>");
                }
                var ink = $(this).find(".ink");
                ink.removeClass("animate-ink");
                if (!ink.height() && !ink.width()) {
                    var d = Math.max($(this).outerWidth(), $(this).outerHeight());
                    ink.css({height: d, width: d});
                }
                var x = e.pageX - $(this).offset().left - ink.width() / 2;
                var y = e.pageY - $(this).offset().top - ink.height() / 2;
                ink.css({top: y + 'px', left: x + 'px'}).addClass("animate-ink");
            })
        },
        filterList: function($menu) {
            if ($menu == null) $menu = $(this.element);
            var $body = $menu.find(".jquery-accordion-menu-body");
            $body.find("li").click(function(){
                $body.find("li.active").removeClass("active")
                $(this).addClass("active");
            })
            var $header = $menu.find(".jquery-accordion-menu-header");
            var form = $("<form>").attr({"class":"filter-form", "action":"#"}),
                input = $("<input>").attr({"class":"filter-input", "type":"text"});

            $(form).append(input).appendTo($header);
            var settings = this.settings;
            input.on('change keyup', function() {
                $body.find('.filtered').removeClass('filtered');
                var filter = $(this).val();
                if (filter) {
                    var $matches = $body.find("a:accordionMenuContains('" + filter + "')").parentsUntil($body, "li");
                    //$matches.filter(":not(.active)").addClass("active");
                    $matches.each(function() {
                        $(this).addClass("filtered");
                        if(settings.expandAllFound) $(this).parent('.submenu').show();
                    });
                    $('li', $body).filter(":not(.filtered)").slideUp();
                    $matches.slideDown();
                } else {
                    $body.find("li").slideDown();
                }
                return false;
            });
        }
    });
    $.fn[pluginName] = function(options) {
        this.each(function() {
            if (!$.data(this, "plugin_" + pluginName)) {
                $.data(this, "plugin_" + pluginName, new Plugin(this, options))
            }
        });
        return this
    };
    $.expr[":"].accordionMenuContains = function(a, i, m) {
        return (a.textContent || a.innerText || "").toUpperCase().indexOf(m[3].toUpperCase()) >= 0;
    };
})(jQuery, window, document);