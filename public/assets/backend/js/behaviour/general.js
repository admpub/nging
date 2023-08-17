var App = function () {

	var config = {//Basic Config
		tooltip: true,
		popover: true,
		nanoScroller: true,
		nestableLists: true,
		hiddenElements: true,
		bootstrapSwitch: true,
		dateTime: true,
		select2: true,
		tags: true,
		slider: true
	};

	/*Form Wizard*/
	var wizard = function () {
		//Fuel UX
		$('.wizard-ux').wizard();

		$('.wizard-ux').on('changed', function () {
			//delete $.fn.slider;
			$('.bslider').slider();
		});

		$(".wizard-next").on('click',function (e) {
			$('.wizard-ux').wizard('next');
			e.preventDefault();
		});

		$(".wizard-previous").on('click',function (e) {
			$('.wizard-ux').wizard('previous');
			e.preventDefault();
		});
	};//End of wizard

	function toggleSideBar(_this) {
		var b = $("#sidebar-collapse")[0];
		var w = $("#cl-wrapper");
		if (w.hasClass("sb-collapsed")) {
			$(".fa", b).addClass("fa-angle-left").removeClass("fa-angle-right");
			w.removeClass("sb-collapsed");
		} else {
			$(".fa", b).removeClass("fa-angle-left").addClass("fa-angle-right");
			w.addClass("sb-collapsed");
		}
	}

	function pageAside() {
		var pageKey = window.location.pathname.replace(/^\//, '').replace(/[^\w]/g, '-');// + '-' + window.location.search.replace(/[^\w]/g, '_');
		$('.page-aside').each(function (index) {
			var pKey = $(this).attr('page-key') || pageKey;
			var aside = $(this), key = pKey + '.page-aside-' + index;
			$(this).find('.header > .collapse-button').on('click', function () {
				aside.addClass('collapsed');
				store.set(key, 'collapsed');
				aside.trigger('collapsed');
				tableReponsive();
			});
			$(this).children('.collapsed-button').on('click', function () {
				aside.removeClass('collapsed');
				store.set(key, '');
				aside.trigger('expanded');
				tableReponsive();
			});
			if (store.get(key) == 'collapsed') {
				$(this).find('.header > .collapse-button').trigger('click');
			}
		});
	}

	function tableReponsiveInit() {
		var aside = $('#main-container > .page-aside');
		var asideWidth = aside.width();
		//if($('#pcont').width()>bodyWidth){
		if (asideWidth < 270) {
			aside.find('.header > .collapse-button').trigger('click');
			return;
		}
		tableReponsive();
	}

	function firstChildrenTable(jqObj) {
		var table = jqObj.children();
		if(table.length<1) return null;
		switch (table[0].tagName.toUpperCase()) {
			case 'TABLE': return table;
			default: return firstChildrenTable(table)
		}
	}

	function tableReponsive(elem) {
		if (elem == null) elem = '#cl-wrapper > .cl-body .table-responsive';
		var jqObj = (elem instanceof jQuery) ? elem : $(elem);
		var pcontLeft = $('#pcont').offset().left, contWidth = $(window).width() - pcontLeft;
		jqObj.each(function () {
			var marginWidth = $(this).offset().left - pcontLeft;
			var bodyWidth = contWidth - marginWidth * 2;
			var oldWidth = $(this).data('width');
			if (oldWidth && oldWidth == bodyWidth) return;
			$(this).data('width', bodyWidth);
			var table = firstChildrenTable($(this));
			if (table == null) return;
			if (table.outerWidth() > bodyWidth) {
				$(this).addClass('overflow').css('max-width', bodyWidth);
			} else {
				$(this).removeClass('overflow').css('max-width', bodyWidth);
			}
		});
	}

	/*SubMenu hover */
	var tool = $("<div id='sub-menu-nav' style='position:fixed;z-index:9999;'></div>");

	function showMenu(_this, e) {
		if (($("#cl-wrapper").hasClass("sb-collapsed") || ($(window).width() > 755 && $(window).width() < 963)) && $("ul", _this).length > 0) {
			$(_this).removeClass("ocult");
			var menu = $("ul", _this);
			if (!$(".dropdown-header", _this).length) {
				var head = '<li class="dropdown-header">' + $(_this).children().html() + "</li>";
				menu.prepend(head);
			}

			tool.appendTo("body");
			var top = ($(_this).offset().top + 8) - $(window).scrollTop();
			var left = $(_this).width();

			tool.css({
				'top': top,
				'left': left + 8
			});
			tool.html('<ul class="sub-menu">' + menu.html() + '</ul>');
			tool.show();

			menu.css('top', top);
		} else {
			tool.hide();
		}
	}

	function returnToTopButton() {
		/*Return to top*/
		var offset = 220;
		var duration = 500;
		var button = $('<a href="#" class="back-to-top"><i class="fa fa-angle-up"></i></a>');
		button.appendTo("body");

		$(window).on('scroll', function () {
			//console.log($(this).scrollTop(), offset)
			if ($(this).scrollTop() > offset) {
				$('.back-to-top').fadeIn(duration);
			} else {
				$('.back-to-top').fadeOut(duration);
			}
		});

		$('.back-to-top').on('click', function (event) {
			event.preventDefault();
			$('html, body').animate({ scrollTop: 0 }, duration);
			return false;
		});
	}

	var cachedLang = null, previousPlotPoint = null; 
	return {
		clientID: {},
		i18n: {
			SYS_INFO: 'System Information', 
			UPLOAD_ERR: 'Upload Error', 
			PLEASE_SELECT_FOR_OPERATE: 'Please select the item you want to operate', 
			PLEASE_SELECT_FOR_REMOVE: 'Please select the item you want to delete', 
			CONFIRM_REMOVE: 'Are you sure you want to delete them?', 
			SELECTED_ITEMS: 'You have selected %d items', 
			SUCCESS: 'The operation was successful', 
			FAILURE: 'Operation failed', 
			UPLOADING:'File uploading, please wait...', 
			UPLOAD_SUCCEED:'Upload successfully', 
			BUTTON_UPLOAD:'Upload' 
		},
		lang: 'en',
		sprintf: sprintfWrapper.init,
		t: function (key) {
			if (typeof (App.i18n[key]) == 'undefined') {
				if (arguments.length < 2) return key;
				return App.sprintf.apply(this, arguments);
			}
			if (arguments.length < 2) return App.i18n[key];
			arguments[0] = App.i18n[key];
			return App.sprintf.apply(this, arguments);
		},
		langInfo: function () {
			if (cachedLang != null) return cachedLang;
			var _lang = App.lang.split('-', 2);
			cachedLang = { encoding: _lang[0], country: '' };
			if (_lang.length > 1) cachedLang.country = _lang[1].toUpperCase();
			return cachedLang;
		},
		langTag: function (seperator) {
			var l = App.langInfo();
			if (l.country) {
				if (seperator == null) {
					seperator = '-';
				}
				return l.encoding + seperator + l.country;
			}
			return l.encoding;
		},
		initTool: function () {
			tool.on("mouseenter", function (e) {
				$(this).addClass("over");
			}).on("mouseleave", function () {
				$(this).removeClass("over");
				tool.fadeOut("fast");
			});
			$(document).on('click',function () {
				tool.hide();
			});
			$(document).on('touchstart click', function (e) {
				tool.fadeOut("fast");
			});
			tool.on('click',function (e) {
				e.stopPropagation();
			});
		},
		initLeftNav: function () {
			/*VERTICAL MENU*/
			$(".cl-vnavigation li ul").each(function () {
				$(this).parent().addClass("parent");
			});

			$(".cl-vnavigation li ul li.active").each(function () {
				$(this).parent().show().parent().addClass("open");
				//setTimeout(function(){updateHeight();},200);
			});
			if (!$(".cl-vnavigation").data('initclick')) {
				$(".cl-vnavigation").data('initclick', true);
				$(".cl-vnavigation").on("click", ".parent > a", function (e) {
					$(".cl-vnavigation .parent.open > ul").not($(this).parent().find("ul")).slideUp(300, 'swing', function () {
						$(this).parent().removeClass("open");
					});

					var ul = $(this).parent().find("ul");
					ul.slideToggle(300, 'swing', function () {
						var p = $(this).parent();
						if (p.hasClass("open")) {
							p.removeClass("open");
						} else {
							p.addClass("open");
						}
						//var menuH = $("#cl-wrapper .menu-space .content").height();
						// var height = ($(document).height() < $(window).height())?$(window).height():menuH;
						//updateHeight();
						$("#cl-wrapper .nscroller").nanoScroller({ preventPageScrolling: true });
					});
					e.preventDefault();
				});
			}
			$(".cl-vnavigation li").on("mouseenter", function (e) {
				showMenu(this, e);
			}).on("mouseleave", function (e) {
				tool.removeClass("over");
				setTimeout(function () {
					if (!tool.hasClass("over") && !$(".cl-vnavigation li:hover").length > 0) {
						tool.hide();
					}
				}, 500);
			});

			$(".cl-vnavigation li").on('click',function (e) {
				if ((($("#cl-wrapper").hasClass("sb-collapsed") || ($(window).width() > 755 && $(window).width() < 963)) && $("ul", this).length > 0) && !($(window).width() < 755)) {
					showMenu(this, e);
					e.stopPropagation();
				}
			});
		},
		initLeftNavAjax: function (activeURL, elem) {
			App.markNavByURL(activeURL);
			App.attachPjax(elem, {
				onclick: function (obj) {
					//console.log($(obj).data('marknav'))
					if ($(obj).data('marknav')) {
						App.unmarkNav($(obj), $(obj).data('marknav'));
						App.markNav($(obj), $(obj).data('marknav'));
					}
				},
				onend: function (evt, xhr, opt) {
					$(opt.container).find('[data-popover="popover"]').popover();
					$(opt.container).find('.ttip, [data-toggle="tooltip"]').tooltip();
				}
			});
			App.attachAjaxURL(elem);
		},
		init: function (options) {
			//Extends basic config with options
			$.extend(config, options);
			App.initLeftNav();
			App.initTool();
			App.showRequriedInputStar();
			/*Small devices toggle*/
			$(".cl-toggle").on('click',function (e) {
				var ul = $(".cl-vnavigation");
				ul.slideToggle(300, 'swing', function () { });
				e.preventDefault();
			});

			/*Collapse sidebar*/
			$("#sidebar-collapse").on('click',function () {
				toggleSideBar();
			});

			if ($("#cl-wrapper").hasClass("fixed-menu")) {
				var scroll = $("#cl-wrapper .menu-space");
				scroll.addClass("nano nscroller");

				function updateHeight() {
					var button = $("#cl-wrapper .collapse-button");
					var collapseH = button.outerHeight();
					var navH = $("#head-nav").height();
					var height = $(window).height() - ((button.is(":visible")) ? collapseH : 0) - navH;
					scroll.css("height", height);
					$("#cl-wrapper .nscroller").nanoScroller({ preventPageScrolling: true });
				}

				$(window).on('resize',function () {
					updateHeight();
				});

				updateHeight();
				$("#cl-wrapper .nscroller").nanoScroller({ preventPageScrolling: true });
			}else{
				$(window).on('resize',function () {
					if($(window).width()>767){
						var navH = $("#head-nav").height();
						$('#cl-wrapper').css("padding-top", navH);
					}
				});
			}

			returnToTopButton();

			/*Datepicker UI*/
			if ($(".ui-datepicker").length > 0) $(".ui-datepicker").datepicker();

			/*Tooltips*/
			if (config.tooltip) {
				$('.ttip, [data-toggle="tooltip"]').tooltip();
			}

			/*Popover*/
			if (config.popover) {
				$('[data-popover="popover"]').popover();
			}

			/*NanoScroller*/
			if (config.nanoScroller) {
				$(".nscroller:not(.has-scrollbar)").nanoScroller();
			}

			/*Nestable Lists*/
			if (config.nestableLists && $('.dd').length > 0) {
				$('.dd').nestable();
			}

			/*Switch*/
			if (config.bootstrapSwitch) {
				if ($('.switch:not(.has-switch)').length > 0) $('.switch:not(.has-switch)').bootstrapSwitch();
			}

			/*DateTime Picker*/
			if (config.dateTime) {
				if ($(".datetime").length > 0) $(".datetime").datetimepicker({ autoclose: true });
			}

			/*Select2*/
			if (config.select2) {
				if ($(".select2").length > 0) $(".select2").select2({
					width: '100%'
				});
			}

			/*Tags*/
			if (config.tags) {
				if ($(".tags").length > 0) $(".tags").select2({ tags: 0, width: '100%' });
			}

			/*Slider*/
			if (config.slider) {
				if ($('.bslider').length > 0) $('.bslider').slider();
			}

			/*Bind plugins on hidden elements*/
			if (config.hiddenElements) {
				/*Dropdown shown event*/
				$('.dropdown').on('shown.bs.dropdown', function () {
					$(".nscroller").nanoScroller();
				});

				/*Tabs refresh hidden elements*/
				$('.nav-tabs').on('shown.bs.tab', function (e) {
					$(".nscroller").nanoScroller();
				});
			}
			App.autoFixedThead();
			$(window).trigger('scroll');
		},
		autoFixedThead: function (prefix) {
			if (prefix == null) prefix = '';
			App.topFloatThead(prefix + 'thead.auto-fixed', $('#head-nav').height());
		},
		pageAside: function (options) {
			pageAside(options);
		},
		tableReponsiveInit: function (options) {
			tableReponsiveInit(options);
		},
		tableReponsive: function (options) {
			tableReponsive(options);
		},
		toggleSideBar: function () {
			toggleSideBar();
		},
		wizard: function () {
			wizard();
		},
		markNavByURL: function (url) {
			if (!url) url = window.location.pathname;
			if (url == '/index') return;
			var leftAnchor=$('#leftnav a[href="' + BACKEND_URL+url + '"]');
			if (leftAnchor.length<1) leftAnchor=$('#leftnav a[href="' + url + '"]');
			App.markNav(leftAnchor, 'left');
			var topAnchor=$('#topnav a[href="' + BACKEND_URL+url + '"]:first');
			if (topAnchor.length<1) topAnchor=$('#topnav a[href="' + url + '"]:first');
			App.markNav(topAnchor, 'top');
		},
		markNav: function (curNavA, position) {
			if (curNavA.length < 1) return;
			var li = curNavA.parent('li').addClass('active');
			switch (position) {
				case 'left':
					li.parent('.from-left').show().parent('li').addClass("open");
					var project = $('#leftnav').attr('data-project');
					//console.log(project);
					var activeProject = $('#topnav li[data-project="' + project + '"]');
					if (activeProject.length > 0 && !activeProject.hasClass('active')) {
						activeProject.addClass('active').siblings('.active').removeClass('active').find('.active').removeClass('active');
					}
					break;

				case 'top':
					li.parent('.from-top').parent('li').addClass("active").siblings('li.active').removeClass('active');
					break;
			}
		},
		unmarkNav: function (curNavA, position) {
			var li = curNavA.parent('li');
			var siblings = li.siblings('li.active');
			if (siblings.length > 0) {
				siblings.removeClass('active');
				return;
			}
			switch (position) {
				case 'left':
					// 点击的左侧边栏菜单
					if (li.parent('ul.sub-menu').length > 0) {
						var op2 = $('.col-menu-2').children('li.open');
						if (op2.length > 0) op2.removeClass('open').find('li.active').removeClass('active');
					}
					$('#leftnav > .open').removeClass('open').children('ul.sub-menu').hide().children('li.active').removeClass('active');
					$('#leftnav .active').removeClass('active');
					break;

				case 'top':
					var topnavDropdown = li.parent('ul.dropdown-menu.from-top');
					if (topnavDropdown.length > 0) {
						siblings = topnavDropdown.parent('li').addClass('active').siblings('li.active');
					}
					if (siblings.length > 0) {
						siblings.removeClass('active').children('ul.dropdown-menu.from-top').children('li.active').removeClass('active');
					}
					break;
			}
		},
		message: function (options, sticky) {
			if (typeof (options) == 'string') {
				switch(options){
					case 'remove':
						var number = sticky;
						return $.gritter.remove(number);
					case 'clear': //$.gritter.removeAll({before_close:function(wrap){},after_close:function(){}});
						return $.gritter.removeAll(sticky||{});
				}
			}
			var defaults = {
				title: App.i18n.SYS_INFO,
				text: '',
				image: '',
				class_name: 'clean',//primary|info|danger|warning|success|dark
				sticky: false // 是否保持显示(不自动关闭)
				//,time: 1000,speed: 500,position: 'bottom-right'
			};
			if (typeof (options) != "object") options = { text: options, class_name: 'clean'};
			if (typeof (options.type) != "undefined" && options.type) options.class_name = options.type;
			options = $.extend({}, defaults, options || {});
			switch (options.class_name) {
				case 'dark':
				case 'primary':
				case 'clean':
				case 'info':
					if (options.title) options.title = '<i class="fa fa-info-circle"></i> ' + App.t(options.title); break;

				case 'error':
					options.class_name = 'danger';
				case 'danger':
					if (options.title) options.title = '<i class="fa fa-comment-o"></i> ' + App.t(options.title); break;

				case 'warning':
					if (options.title) options.title = '<i class="fa fa-warning"></i> ' + App.t(options.title); break;

				case 'success':
					if (options.title) options.title = '<i class="fa fa-check"></i> ' + App.t(options.title); break;
			}
			if (sticky != null) options.sticky = sticky;
			if (options.text) options.text = App.t(options.text);
			var number = $.gritter.add(options);
			return number;
		},
		attachAjaxURL: function (elem) {
			if (elem == null) elem = document;
			$(elem).on('click', '[data-ajax-url]', function () {
				var a = $(this), confirmMsg = a.data('ajax-confirm');
				if(a.data('processing')){
					alert(App.t('Processing, please wait for the operation to complete'));
					return;
				}
				if(confirmMsg && !confirm(confirmMsg)) return;
				a.data('processing',true);
				var url = a.data('ajax-url'), method = a.data('ajax-method') || 'get', params = a.data('ajax-params') || {}, title = a.attr('title')||App.i18n.SYS_INFO, accept = a.data('ajax-accept') || 'html', target = a.data('ajax-target'), callback = a.data('ajax-callback'), toggle = a.data('ajax-toggle'), onsuccess = a.data('ajax-onsuccess');
				if (!title) title = a.text();
				var fa = a.children('.fa');
				var hasIcon = toggle && fa.length>0;
				if (hasIcon){
					fa.addClass('fa-spin')
				}else{
					App.loading('show');
				}
				a.trigger('processing');
				if (typeof params === "function") params = params.call(this, arguments);
				$[method](url, params || {}, function (r) {
					a.data('processing',false);
					a.trigger('finished',arguments);
					if (hasIcon){
						fa.removeClass('fa-spin');
					}else{
						App.loading('hide');
					}
					if (callback) return callback.call(this, arguments);
					if (target) {
						var data;
						if (accept == 'json') {
							if (r.Code != 1) {
								return App.message({ title: title, text: r.Info, type: 'error', time: 5000, sticky: false });
							}
							data = r.Data;
						} else {
							data = r;
						}
						$(target).html(data);
						a.trigger('partial.loaded', arguments);
						$(target).trigger('partial.loaded', arguments);
						if(onsuccess) window.setTimeout(onsuccess,0);
						return;
					}
					if(onsuccess) window.setTimeout(onsuccess,3000);
					if (accept == 'json') {
						return App.message({ title: title, text: r.Info, type: r.Code == 1 ? 'success' : 'error', time: 5000, sticky: false });
					}
					App.message({ title: title, text: r, time: 5000, sticky: false });
				}, accept).error(function (xhr, status, info) {
					a.data('processing',false);
					a.trigger('finished',arguments);
					if (hasIcon){
						fa.removeClass('fa-spin');
					}else{
						App.loading('hide');
					}
					App.message({ title: title, text: xhr.responseText, type: 'error', time: 5000, sticky: false });
				});
			});
		},
		attachPjax: function (elem, callbacks, timeout) {
			if (!$.support.pjax) return;
			if (elem == null) elem = 'a';
			if (timeout == null) timeout = 5000;
			var defaults = { onclick: null, onsend: null, oncomplete: null, ontimeout: null, onstart: null, onend: null };
			var options = $.extend({}, defaults, callbacks || {});
			$(document).on('click', elem + '[data-pjax]', function (event) {
				var container = $(this).data('pjax'), keepjs = $(this).data('keepjs');
				var onclick = $(this).data('onclick'), toggleClass = $(this).data('toggleclass');
				$.pjax.click(event, container, { timeout: timeout, keepjs: keepjs });
				if (options.onclick) options.onclick(this);
				if (onclick && typeof (window[onclick]) == 'function') window[onclick](this);
				if (toggleClass) {
					var arr = toggleClass.split(':'),parent,target;
					if (arr.length>1) {
						parent = arr[0];
						toggleClass = arr[1];
					}
					toggleClass = toggleClass.replace(/^\./g,'');
					if(parent){
						var index = parent.indexOf('.'),parentElem;
						if (index > 0) {
							parentElem = parent.substring(index+1);
							parent = parent.substring(0,index);
						}
						switch(parent){
							case 'parent':target=$(this).parent(parentElem);break;
							case 'parents':target=$(this).parents(parentElem);break;
							case 'closest':target=$(this).closest(parentElem);break;
							default:target=$(this).parent(parentElem);break;
						}
					}else{
						target = $(this);
					}
					target.addClass(toggleClass).siblings('.'+toggleClass).removeClass(toggleClass);
				}
				$('.sp_result_area').remove();
				$('.tox').remove();
				$('.select2-hidden-accessible').remove();
				$('.select2-sizer').remove();
				$('.select2-drop').remove();
				$('#select2-drop-mask').remove();
			}).on('pjax:send', function (evt, xhr, option) {
				App.loading('show');
				if (options.onsend) options.onsend(evt, xhr, option);
			}).on('pjax:complete', function (evt, xhr, textStatus, option) {
				App.loading('hide');
				if (options.oncomplete) options.oncomplete(evt, xhr, textStatus, option);
			}).on('pjax:timeout', function (evt, xhr, option) {
				console.log('timeout');
				App.loading('hide');
				if (options.ontimeout) options.ontimeout(evt, xhr, option);
			}).on('pjax:start', function (evt, xhr, option) {
				if (options.onstart) options.onstart(evt, xhr, option);
			}).on('pjax:end', function (evt, xhr, option) {
				App.loading('hide');
				if (options.onend) options.onend(evt, xhr, option);
				if (option.container) {
					App.bottomFloat(option.container + ' .pagination');
					App.bottomFloat(option.container + ' .form-submit-group', 0, true);
					$(option.container + ' .switch:not(.has-switch)').bootstrapSwitch();
					App.autoFixedThead(option.container + ' ');
				}
				if (option.type == 'GET') $('#global-search-form').attr('action', option.url);
			});
		},
		wsURL: function (url) {
			var protocol = 'ws:';
			if (window.location.protocol == 'https:') protocol = 'wss:';
			var p = String(url).indexOf('//');
			if (p == -1) {
				url = protocol + "//" + window.location.host + url;
			} else {
				url = protocol + String(url).substring(p);
			}
			return url;
		},
		websocket: function (showmsg, url, onopen) {
			url = App.wsURL(url);
			var ws = new WebSocket(url);
			ws.onopen = function (evt) {
				console.log('Websocket Server is connected');
				if (onopen != null && typeof onopen === "function") onopen.apply(this, arguments);
			};
			ws.onclose = function (evt) {
				console.log('Websocket Server is disconnected');
			};
			ws.onmessage = function (evt) {
				showmsg(evt.data);
			};
			ws.onerror = function (evt) {
				console.dir(evt);
			};
			if (onopen != null && typeof (onopen) == 'object') {
				ws = $.extend({}, ws, onopen);
			}
			return ws;
		},
		notifyListen: function () {
			var messageCount = {notify: 0, element: 0, modal: 0},  
			messageMax = {notify: 20, element: 50, modal: 50};
			App.websocket(function (message) {
				//console.dir(message);
				var m = $.parseJSON(message);
				if (!m) {
					App.message({text:message||'Websocket Server is disconnected',type:'error'});
					return false;
				}
				if (typeof(App.clientID['notify']) == 'undefined') {
					App.clientID['notify'] = m.client_id;
				}
				if (typeof(m.content) == 'undefined' || !m.content) {
					return false;
				}
				switch (m.mode) {
					case '-':
						break;
					case 'element':
						var c = $('#notify-element-' + m.type);
						if (c.length < 1) {
							var callback = 'recv_notice_' + m.type;
							if (typeof (window[callback]) != 'undefined') {
								return window[callback](m);
							}
							if (m.status > 0) {
								console.info(m.content);
							} else {
								console.error(m.content);
							}
							return false;
						}
						if (messageCount[m.mode] >= messageMax[m.mode]) {
							c.find('li:first').remove();
						}
						if (m.title) {
							var badge = 'badge-danger';
							if (m.status > 0) badge = 'badge-success';
							message = '<span class="badge ' + badge + '">' + App.text2html(m.title) + '</span> ' + App.text2html(m.content);
						} else {
							message = App.text2html(m.content);
						}
						c.append('<li>' + message + '</li>');
						messageCount[m.mode]++;
						break;
					case 'modal':
						var c = $('#notify-modal-' + m.type);
						if (c.length < 1) {
							var callback = 'recv_notice_' + m.type;
							if (typeof (window[callback]) != 'undefined') {
								return window[callback](m);
							}
							if (m.status > 0) {
								console.info(m.content);
							} else {
								console.error(m.content);
							}
							return false;
						}
						if (m.title) {
							var badge = 'badge-danger';
							if (m.status > 0) badge = 'badge-success';
							message = '<span class="badge ' + badge + '">' + App.text2html(m.title) + '</span> ' + App.text2html(m.content);
						} else {
							message = App.text2html(m.content);
						}
						if (!c.data('shown')) {
							messageCount[m.mode] = 0;
							c.data('shown', true);
							var mbody = c.find('.modal-body'), mbodyUL = mbody.children('ul.modal-body-ul');
							if (mbodyUL.length < 1) {
								mbody.html('<ul class="modal-body-ul" id="notify-modal-' + m.type + '-container"><li>' + message + '</li></ul>');
							} else {
								mbodyUL.html('<li>' + message + '</li>');
							}
							c.niftyModal('show', {
								afterOpen: function (modal) { },
								afterClose: function (modal) {
									c.data('shown', false);
								}
							});
						} else {
							var cc = $('#notify-modal-' + m.type + '-container');
							if (messageCount[m.mode] >= messageMax[m.mode]) {
								cc.find('li:first').remove();
							}
							cc.append('<li>' + message + '</li>');
						}
						messageCount[m.mode]++;
						break;
					case 'notify':
					default:
						if ('notify' != m.mode) m.mode = 'notify';
						var c = $('#notice-message-container');
						if (c.length < 1) {
							App.message({ title: App.i18n.SYS_INFO, text: '<ul id="notice-message-container" class="no-list-style" style="max-height:500px;overflow-y:auto;overflow-x:hidden"></ul>', sticky: true });
							c = $('#notice-message-container');
						}
						if (messageCount[m.mode] >= messageMax[m.mode]) {
							c.find('li:first').remove();
						}
						if (m.title) {
							var badge = 'badge-danger';
							if (m.status > 0) badge = 'badge-success';
							message = '<span class="badge ' + badge + '">' + App.text2html(m.title) + '</span>' + App.text2html(m.content);
						} else {
							message = App.text2html(m.content);
						}
						c.append('<li>' + message + '</li>');
						messageCount[m.mode]++;
						break;
				}
				return true;
			}, BACKEND_URL + '/user/notice');
		},
		text2html: function (text, noescape) {
			text = String(text);
			if(!noescape) text = text.replace(/</g, '&lt;').replace(/>/g, '&gt;');
			return App.textNl2br(text);
		},
		ifTextNl2br: function (text) {
			text = String(text);
			if (/<[^>]+>/.test(text)) return text;
			return App.textNl2br(text);
		},
		textNl2br: function (text) {
			return text.replace(/\n/g, '<br />').replace(/  /g, '&nbsp; ').replace(/\t/g, '&nbsp; &nbsp; ');
		},
		trimSpace: function (text) {
			return String(text).replace(/^[\s]+|[\s]+$/g,'');
		},
		checkedAll: function (ctrl, target) {
			return $(target).not(':disabled').prop('checked', $(ctrl).prop('checked'));
		},
		attachCheckedAll: function (ctrl, target, showNumElem) {
			$(ctrl).on('ifChecked ifUnchecked click', function () {
				App.checkedAll(this, target);
				if (showNumElem) $(showNumElem).text($(target + ':checked').length);
			});
		},
		alertBlock: function (content, title, type) {
			switch (type) {
				case 'info':
					if (title == null) title = 'Info!';
					return '<div class="alert alert-info">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<i class="fa fa-info-circle sign"></i><strong>'+ title + '</strong> ' + content + '</div>';
				case 'warn':
					if (title == null) title = 'Alert!';
					return '<div class="alert alert-warning">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<i class="fa fa-warning sign"></i><strong>'+ title + '</strong> ' + content + '</div>';
				case 'error':
					if (title == null) title = 'Error!';
					return '<div class="alert alert-danger">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<i class="fa fa-times-circle sign"></i><strong>'+ title + '</strong> ' + content + '</div>';
				default:
					if (title == null) title = 'Success!';
					return '<div class="alert alert-success">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<i class="fa fa-check sign"></i><strong>'+ title + '</strong> ' + content + '</div>';
			}
		},
		alertBlockx: function (content, title, type) {
			switch (type) {
				case 'info':
					if (title == null) title = 'Info!';
					return '<div class="alert alert-info alert-white rounded">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<div class="icon"><i class="fa fa-info-circle"></i></div>\
								<strong>'+ title + '</strong> ' + content + '</div>';
				case 'warn':
					if (title == null) title = 'Alert!';
					return '<div class="alert alert-warning alert-white rounded">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<div class="icon"><i class="fa fa-warning"></i></div>\
								<strong>'+ title + '</strong> ' + content + '</div>';
				case 'error':
					if (title == null) title = 'Error!';
					return '<div class="alert alert-danger alert-white rounded">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<div class="icon"><i class="fa fa-times-circle"></i></div>\
								<strong>'+ title + '</strong> ' + content + '</div>';
				default:
					if (title == null) title = 'Success!';
					return '<div class="alert alert-success alert-white rounded">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<div class="icon"><i class="fa fa-check"></i></div>\
								<strong>'+ title + '</strong> ' + content + '</div>';
			}
		},
		showDbLog: function (results, container) {
			if (container == null) container = '.block-flat:first';
			if (typeof (results.length) == 'undefined') results = [results];
			for (var j = 0; j < results.length; j++) {
				var result = results[j];
				var s = result.Started;
				if (result.SQLs && result.SQLs.length > 0) {
					for (var i = 0; i < result.SQLs.length; i++) s += '<code class="wrap">' + result.SQLs[i] + '</code>';
				} else {
					s += '<code class="wrap">' + result.SQL + '</code>';
				}
				var t = 'success';
				if (result.Error) {
					s += '(' + result.Error + ')';
					t = 'error'
				} else {
					s += '(' + result.Elapsed + ')';
				}
				$(container).before(App.alertBlockx(s, null, t));
			}
		},
		loading: function (op) {
			var obj = $('#loading-status');
			switch (op) {
				case 'show':
					if (obj.length > 0) {
						obj.show();
					} else {
						$('body').append('<div id="loading-status"><i class="fa fa-spinner fa-spin fa-3x"></i></div>');
					}
					break;
				case 'hide':
					if (obj.length > 0) {
						obj.hide();
					}
			}
		},
		insertAtCursor: function (myField, myValue, posStart, posEnd) {
			if (typeof TextAreaEditor != 'undefined') {
				TextAreaEditor.setSelectText(myField, myValue, posStart, posEnd);
				return;
			}
			/* IE support */
			if (document.selection) {
				myField.focus();
				sel = document.selection.createRange();
				sel.text = myValue;
				sel.select();
			} /* MOZILLA/NETSCAPE support */
			else if (myField.selectionStart || myField.selectionStart == '0') {
				var startPos = myField.selectionStart;
				var endPos = myField.selectionEnd; /* save scrollTop before insert */
				var restoreTop = myField.scrollTop;
				myField.value = myField.value.substring(0, startPos) + myValue + myField.value.substring(endPos, myField.value.length);
				if (restoreTop > 0) myField.scrollTop = restoreTop;
				myField.focus();
				myField.selectionStart = startPos + myValue.length;
				myField.selectionEnd = startPos + myValue.length;
			} else {
				myField.value += myValue;
				myField.focus();
			}
		},
		searchFS: function (elem, size, type, url, before) {
			if (size == null) size = 10;
			if (url == null) url = BACKEND_URL + '/user/autocomplete_path';
			$(elem).typeahead({
				hint: true, highlight: true, minLength: 1
			}, {
				source: function (query, sync, async) {
					var data = { query: query, size: size, type: type };
					$.ajax({
						url: url,
						type: 'get',
						data: before ? before(data) : data,
						dataType: 'json',
						async: false,
						success: function (data) {
							var arr = [];
							if (!data.Data) return;
							$.each(data.Data, function (index, val) {
								arr.push(val);
							});
							sync(arr);
						}
					});
				}, limit: size
			});
		},
		randomString: function (len) {
			len = len || 32;
			var $chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678'; //默认去掉了容易混淆的字符oOLl,9gq,Vv,Uu,I1
			var maxPos = $chars.length;
			var pwd = '';
			for (i = 0; i < len; i++) {
				pwd += $chars.charAt(Math.floor(Math.random() * maxPos));
			}
			return pwd;
		},
		bottomFloat: function (elems, top, autoWith) {
			if ($(elems).length < 1) return;
			if (top == null) top = 0;
			$(elems).not('[disabled-fixed]').each(function () {
				$(this).attr('disabled-fixed', 'fixed');
				var elem = this;
				var _offset = $(elem).height() + top;
				var offsetY = $(elem).offset().top + _offset;
				var w = $(elem).outerWidth(), h = $(elem).outerHeight();
				if (!autoWith) autoWith = $(elem).data('auto-width');
				if (autoWith) $(elem).css('width', w);
				$(window).on('scroll', function () {
					var scrollH = $(this).scrollTop() + $(window).height();
					if (scrollH >= offsetY) {
						if ($(elem).hasClass('always-bottom')) {
							$(elem).removeClass('always-bottom');
							$(elem).next('.fixed-placeholder').hide();
						}
						return;
					}
					if (!$(elem).hasClass('always-bottom')) {
						$(elem).addClass('always-bottom');
						if ($(elem).next('.fixed-placeholder').length > 0) {
							$(elem).next('.fixed-placeholder').show();
						} else {
							$(elem).after('<div style="width:' + w + 'px;height:' + h + 'px" class="fixed-placeholder"></div>');
						}
					}
				});//on-scroll
			});//end-each
			//$(window).trigger('scroll');
		},
		topFloat: function (elems, top, autoWith) {
			if ($(elems).length < 1) return;
			if (top == null) top = 0;
			$(elems).not('[disabled-fixed]').each(function () {
				$(this).attr('disabled-fixed', 'fixed');
				var elem = this;
				var offsetY = $(elem).offset().top;
				var w = $(elem).outerWidth(), h = $(elem).outerHeight();
				if (!autoWith) autoWith = $(elem).data('auto-width');
				if (autoWith) $(elem).css('width', w);
				$(window).on('scroll', function () {
					var scrollH = $(this).scrollTop();
					if (scrollH <= offsetY) {
						if ($(elem).hasClass('always-top')) {
							$(elem).removeClass('always-top');
							$(elem).next('.fixed-placeholder').hide();
						}
						return;
					}
					if (!$(elem).hasClass('always-top')) {
						$(elem).addClass('always-top').css('top', top);
						if ($(elem).next('.fixed-placeholder').length > 0) {
							$(elem).next('.fixed-placeholder').show();
						} else {
							$(elem).after('<div style="width:' + w + 'px;height:' + h + 'px" class="fixed-placeholder"></div>');
						}
					}
				});//on-scroll
			});//end-each
			//$(window).trigger('scroll');
		},
		topFloatRawThead: function (elems, top) {
			if ($(elems).length < 1) return;
			if (top == null) top = 0;
			$(elems).not('[disabled-fixed]').each(function () {
				$(this).attr('disabled-fixed', 'fixed');
				var elem = this, table = $(elem).parent('table'), reponsive = table.parents('.table-responsive');
				var scrollable = reponsive.length > 0;
				var offsetY = $(elem).offset().top, maxOffsetY = table.height() + offsetY - $(elem).outerHeight() * 2;
				$(elem).css({ 'background-color': 'white' });
				var setSize = function (init) {
					if (init == null) init = false;
					if (scrollable) {
						tableReponsive(reponsive);
						if (!reponsive.hasClass('overflow')) {
							$(elem).css({ 'width': 'auto', 'overflow-x': 'unset' });
						}
					}
					var width = $(elem).outerWidth(), ratio = 1;
					if (!init) {
						if (Math.abs(table.data('width') - width) > 1) {//避免抖动
							ratio = width / table.data('width');
						}
						if (table.data('offset-left') != $(table).offset().left) {
						} else if (table.data('scroll-left') != $(window).scrollLeft()) {
							$(elem).css({ 'left': table.offset().left - $(window).scrollLeft() });
							if (ratio == 1) return;
						}
					}
					table.data('width', width);//记录宽度
					table.data('offset-left', $(table).offset().left);//记录左侧偏移
					table.data('scroll-left', $(window).scrollLeft());//记录左侧滚动条
					var cols = table.children('colgroup').length>0?table.children('colgroup').children('col'):table.children('col'), tds = $(elem).find('td,th');
					if (cols.length < 1) {
						var html = '';
						tds.each(function () {
							var w = $(this).outerWidth() * ratio;
							html += '<col style="min-width:' + w + 'px;max-width:auto" />';
							$(this).css({ 'min-width': w, 'max-width': 'auto' });
						});
						table.prepend(html);
						return;
					}
					tds.each(function (index) {
						var col = cols.eq(index);
						var w = $(this).outerWidth() * ratio;
						col.css({ 'width': w });
						$(this).css({ 'width': w });
					});
				}
				setSize(true);
				$(window).on('scroll resize', function () {
					setSize();
					var scrollH = $(this).scrollTop();
					if (scrollH <= offsetY || scrollH >= maxOffsetY) {
						if ($(elem).hasClass('always-top')) {
							$(elem).removeClass('always-top');
						}
						if (scrollable) {
							$(elem).off('scroll').data('scroll-reponsive', false);
							reponsive.off('scroll').data('scroll-thead', false);
						}
						return;
					}
					if (table.height() > $(window).height()) {
						if (!$(elem).hasClass('always-top')) $(elem).addClass('always-top');
						var cssOpts = { 'top': top };
						if (scrollable) {
							cssOpts['width'] = reponsive.outerWidth();
							if (reponsive.hasClass('overflow')) cssOpts['overflow-x'] = 'auto';

							if (!$(elem).data('scroll-reponsive')) {
								$(elem).on('scroll', function () {
									reponsive.scrollLeft($(this).scrollLeft());
								}).data('scroll-reponsive', true);
							}
							if (!reponsive.data('scroll-thead')) {
								reponsive.on('scroll', function () {
									$(elem).scrollLeft($(this).scrollLeft());
								}).data('scroll-thead', true);
							}
						}
						$(elem).css(cssOpts);
					}
				});
			});
			$(window).trigger('scroll');
		},
		topFloatThead: function (elems, top, clone) {
			if (!clone) return App.topFloatRawThead(elems, top);
			if ($(elems).length < 1) return;
			if (top == null) top = 0;
			$(elems).not('[disabled-fixed]').each(function () {
				$(this).attr('disabled-fixed', 'fixed');
				var elem = this, table = $(elem).parent('table');
				var offsetY = $(elem).offset().top, maxOffsetY = table.height() + offsetY - $(elem).outerHeight() * 2, cid = $(elem).data('copy');
				if (cid) {
					$('#tableCopy' + cid).remove();
				} else {
					cid = Math.random();
					$(elem).data('copy', cid);
				}
				var eCopy = $('<table class="' + table.attr('class') + ' always-top" style="background-color:white" id="tableCopy' + cid + '"></table>');
				var hCopy = $(elem).clone();
				eCopy.append(hCopy);
				var setSize = function (init) {
					if (init == null) init = false;
					if (!init) {
						if (eCopy.data('offset-left') != $(elem).offset().left) {
						} else if (eCopy.data('scroll-left') != $(window).scrollLeft()) {
							eCopy.css({ 'left': $(elem).offset().left - $(window).scrollLeft() });
							return;
						} else {
							return;
						}
					}
					eCopy.data('offset-left', $(elem).offset().left);//记录左侧偏移
					eCopy.data('scroll-left', $(window).scrollLeft());//记录左侧滚动条
					var cols = hCopy.find('td,th'), rawCols = $(elem).find('td,th');
					rawCols.each(function (index) {
						var col = cols.eq(index);
						col.css({'width': $(this).outerWidth()});
						if (!init) return;
						var chk = col.find('input:checkbox');
						if (chk.length < 1) return;
						var rawChk = rawCols.find('input:checkbox');
						chk.each(function (idx) {
							rawChk.eq(idx).on('click change', function () {
								chk.prop('checked', $(this).prop('checked'));
							});
						});
					});
					var offsetX = $(elem).offset().left - $(window).scrollLeft();
					var w = $(elem).outerWidth(), h = $(elem).outerHeight()
					eCopy.css({ 'top': top, 'left': offsetX, 'width': w, 'height': h });
				}
				setSize(true);
				eCopy.hide();
				table.after(eCopy);
				$(window).on('scroll', function () {
					setSize();
					var scrollH = $(this).scrollTop();
					if (scrollH <= offsetY || scrollH >= maxOffsetY) return eCopy.hide();
					eCopy.show();
				});
			});
			//$(window).trigger('scroll');
		},
		getImgNaturalDimensions: function (oImg, callback) {
			if (!oImg.naturalWidth) { // 现代浏览器
				callback({ w: oImg.naturalWidth, h: oImg.naturalHeight });
				return;
			}
			// IE6/7/8
			var nImg = new Image();
			nImg.onload = function () {
				callback({ w: nImg.width, h: nImg.height });
			}
			nImg.src = oImg.src;
		},
		reportBug: function (url) {
			$.post(url, { "panic": $('#panic-content').html(), "url": window.location.href }, function (r) { }, 'json');
		},
		replaceURLParam: function (name, value, url) {
			if (url == null) url = window.location.href;
			value = encodeURIComponent(value);
			var pos = String(url).indexOf('?');
			if (pos < 0) return url + '?' + name + '=' + value;
			var q = url.substring(pos), r = new RegExp('([\\?&]' + name + '=)[^&]*(&|$)');
			if (!r.test(q)) return url + '&' + name + '=' + value;
			url = url.substring(0, pos);
			q = q.replace(r, '$1' + value + '$2');
			return url + q;
		},
		switchLang: function (lang) {
			window.location = App.replaceURLParam('lang', lang);
		},
		extends: function (child, parent) {
			//parent.call(this);
			var obj = function () { };
			obj.prototype = parent.prototype;
			child.prototype = new obj();
			child.prototype.constructor = child;
		},
		formatBytes: function (bytes, precision) {
			if (precision == null) precision = 2;
			var units = ["YB", "ZB", "EB", "PB", "TB", "GB", "MB", "KB", "B"];
			var total = units.length;
			for (total--; total > 0 && bytes > 1024.0; total--) {
				bytes /= 1024.0;
			}
			return bytes.toFixed(precision) + units[total];
		},
		format: function (raw, xy, formatters) {
			if (formatters) {
				if (typeof (formatters[xy.y]) == 'function') return formatters[xy.y](raw, xy.x);
				if (typeof (formatters['']) == 'function') return formatters[''](raw, xy.x);
			}
			return raw;
		},
		genTable: function (rows, formatters) {
			var h = '<table class="table table-bordered no-margin">';
			var th = '<thead>', bd = '<tbody>';
			for (var i = 0; i < rows.length; i++) {
				var v = rows[i];
				if (i == 0) {
					for (var k in v) th += '<th class="' + k + '"><strong>' + k + '</strong></th>';
				}
				bd += '<tr>';
				for (var k in v) bd += '<td class="' + k + '">' + App.format(v[k], { 'x': i, 'y': k }, formatters) + '</td>';
				bd += '</tr>';
			}
			th += '</thead>';
			bd += '</tbody>';
			h += th + bd + '</table>';
			return h;
		},
		httpStatusColor: function (code) {
			if (code >= 500) return 'danger';
			if (code >= 400) return 'warning';
			if (code >= 300) return 'info';
			return 'success';
		},
		htmlEncode: function(value){
			return !value ? value : String(value).replace(/&/g, "&amp;").replace(/>/g, "&gt;").replace(/</g, "&lt;").replace(/"/g, "&quot;");
		},
		htmlDecode: function(value){
			return !value ? value : String(value).replace(/&gt;/g, ">").replace(/&lt;/g, "<").replace(/&quot;/g, '"').replace(/&amp;/g, "&");
		},
		logShow: function (elem, trigger, pipe) {
			var title=$(elem).data('modal-title');
			if(title) $('#log-show-modal').find('.modal-header h3').text(title);
			if (!$('#log-show-modal').data('init')) {
				$('#log-show-modal').data('init', true);
				$(window).off().on('resize', function () {
					$('#log-show-modal').css({ height: $(window).height(), width: '100%', 'max-width': '100%', left: 0, top: 0, transform: 'none' });
					$('#log-show-modal').find('.md-content').css('height', $(window).height());
					$('#log-show-content').css('height', $(window).height() - 200);
				});
				$('#log-show-last-lines').on('change', function (r) {
					var target = $(this).data('target');
					if (!target) return;
					var lastLines = $(this).val();
					target.data('last-lines', lastLines);
					target.trigger('click');
				});
				$('#log-show-modal .modal-footer .btn-refresh').on('click', function (r) {
					var target = $('#log-show-last-lines').data('target');
					if (!target) return;
					target.trigger('click');
				});
				$(window).trigger('resize');
			}
			if (pipe == null) pipe = '';
			var done = function (a) {
				var url = $(a).data('url');
				var lastLines = $(a).data('last-lines');
				if (lastLines == null) lastLines = 100;
				$('#log-show-last-lines').data('target', $(a));
				var contentID = 'log-show-content', contentE = '#' + contentID;
				$('#log-show-modal').niftyModal('show', {
					afterOpen: function (modal) {
						$.get(url, { lastLines: lastLines, pipe: pipe }, function (r) {
							if (r.Code == 1) {
								var subTitle = $('#log-show-modal .modal-header .modal-subtitle');
								if (typeof (r.Data.title) != 'undefined') {
									if (r.Data.title) r.Data.title = ' (' + r.Data.title + ')';
									subTitle.html(r.Data.title);
								} else {
									subTitle.empty();
								}
								if (typeof (r.Data.list) != 'undefined') {
									var h = '<div class="table-responsive" id="' + contentID + '">' + App.genTable(r.Data.list, {
										'StatusCode': function (raw, index) {
											return '<span class="label label-' + App.httpStatusColor(raw) + '">' + raw + '</span>';
										},
										'': function (raw, index) {
											return App.htmlEncode(raw);
										}
									}) + '</div>';
									$(contentE).parent('.modal-body').css('padding', 0);
									$(contentE).replaceWith(h);
								} else {
									if ($(contentE)[0].tagName.toUpperCase() != 'TEXTAREA') {
										$(contentE).replaceWith("<textarea name='content' class='form-control' id='" + contentID + "'></textarea>");
									}
									$(contentE).text(r.Data.content);
								}
								$(window).trigger('resize');
								var textarea = $(contentE)[0];
								textarea.scrollTop = textarea.scrollHeight;
							} else {
								$(contentE).text(r.Info);
							}
						}, 'json');
					},
					afterClose: function (modal) { 
						$('#log-show-last-lines').find('option:selected').prop('selected',false);
					}
				});
			};
			if (trigger) return done(elem);
			$(elem).on('click', done(this));
		},
		tableSorting: function (table) {
			table = table == null ? '' : table + ' ';
			function sortAction(sortObj, isDesc) {
				var newCls, oldCls, sortBy;
				if (!isDesc) {
					newCls = 'fa-arrow-up';
					sortBy = 'up';
					oldCls = 'fa-arrow-down';
				} else {
					newCls = 'fa-arrow-down';
					sortBy = 'down';
					oldCls = 'fa-arrow-up';
				}
				if (sortObj.length > 0) {
					var icon = sortObj.children('.fa');
					if (icon.length < 1) {
						sortObj.append('<i class="fa ' + newCls + '"></i>');
					} else {
						icon.removeClass(oldCls).addClass(newCls);
					}
					sortObj.addClass('sort-active sort-' + sortBy);
					sortObj.siblings('.sort-active').removeClass('sort-active').removeClass('sort-up').removeClass('sort-down').find('.fa').remove();
				}
			}
			$(table + '[sort-current!=""]').each(function () {//<thead sort-current="created">
				var current = String($(this).attr('sort-current'));
				var isDesc = current.substring(0, 1) == '-';
				if (isDesc) current = current.substring(1);
				var sortObj = $(this).find('[sort="' + current + '"]');//<th sort="-created">
				if (sortObj.length < 1 && current) {
					sortObj = $(this).find('[sort="-' + current + '"]');
				}
				sortAction(sortObj, isDesc);
			});
			$(table + '[sort-current] [sort]').css('cursor', 'pointer').on('click', function (e) {
				var thead = $(this).parents('[sort-current]');
				var current = thead.attr('sort-current');
				var url = thead.attr('sort-url') || window.location.href;
				var trigger = thead.attr('sort-trigger') || thead.data('sort-trigger');
				var sort = $(this).attr('sort');
				if (current && (current == sort || current == '-' + sort)) {
					var reg = /^\-/;
					current = reg.test(current) ? current.replace(reg, '') : '-' + current;
				} else {
					current = sort;
				}
				thead.attr('sort-current', current);
				url = App.replaceURLParam('sort', current, url);
				if (trigger) {
					thead.data('sort-url', url);
					var isDesc = current.substring(0, 1) == '-';
					sortAction($(this), isDesc);
					window.setTimeout(trigger, 0);
				} else {
					var setto = thead.attr('sort-setto');
					if (setto) {
						var $this=$(this);
						$.get(url,{},function(r,status,xhr){
							if(String(xhr.getResponseHeader('Content-Type')).split(';')[0]=='application/json'){
								try {
									r = JSON.parse(r);
								} catch (error) {
									return App.message({text:error,type:'error'});
								}
								if(r.Code!=1) return App.message({text:r.Info,type:'error'});
								r = r.Data.html;
							}
							$(setto).html(r);
							if($(setto).length>0 && $(setto)[0].tagName.toUpperCase()=='TBODY'){
								var thead = $this.parents('[sort-current]');
								var current = thead.attr('sort-current');
								var isDesc = current.substring(0, 1) == '-';
								sortAction($this, isDesc);
							}
						},'html');
					} else {
						window.location = url;
					}
				}
			});
		},
		resizeModalHeight: function (el) {
			var h = $(window).height() - 200;
			if (h < 200) h = 200;
			var bh = h - 150;
			$(el).css({ "max-height": h + 'px' });
			$(el).find('.modal-body').css({ "max-height": bh + 'px' });
		},
		switchStatus: function (a, type, editURL, callback) {
			if (type == null) type = $(a).data('type');
			var v = $(a).val();
			var checkedValue = $(a).data('v-checked')||v||'N',
			uncheckedValue = $(a).data('v-unchecked')||(checkedValue=='N'?'Y':'N');
			if (type) {
				var tmp=String(type).split('=');//disabled=Y|N
				type=tmp[0];
				if(tmp.length>1&&tmp[1]){
					var optValues=tmp[1].split('|');
					if(optValues[0]==checkedValue){
						checkedValue=optValues[0];
						uncheckedValue=optValues[1];
					}else{
						checkedValue=optValues[1];
						uncheckedValue=optValues[0];
					}
				}
			}
			if (editURL == null) editURL = $(a).data('url');
			var that = $(a), status = a.checked ? checkedValue : uncheckedValue, data = { id: that.data('id') };
			data[type] = status;
			if (String(editURL).charAt(0) != '/') editURL = BACKEND_URL + '/' + editURL;
			$.get(editURL, data, function (r) {
				if (r.Code == 1) {
					that.attr('data-' + type, status);
					that.prop('checked', status == v);
				}
				App.message({ title: App.i18n.SYS_INFO, text: r.Info, time: 5000, sticky: false, class_name: r.Code == 1 ? 'success' : 'error' });
				if (callback) callback.call(a, r);
			}, 'json');
		},
		bindSwitch: function (elem, eventName, editURL, type, callback) {
			if (eventName == null) eventName = 'click';
			var re = new RegExp('switch-([\\w\\d]+)');
			$(elem).on(eventName, function () {
				if (type == null) {
					var matches = String($(this).attr('class')).match(re);
					type = matches[1];
				}
				App.switchStatus(this, type, editURL, callback);
			});
		},
		removeSelected: function (elem, postField, removeURL, callback) {
			return App.opSelected(elem, postField, removeURL, callback, App.i18n.CONFIRM_REMOVE, App.i18n.PLEASE_SELECT_FOR_REMOVE);
		},
		opSelected: function (elem, postField, removeURL, callback, confirmMsg, unselectedMsg) {
			if (removeURL == null) {
				removeURL = window.location.href;
			} else if (String(removeURL).charAt(0) != '/') {
				removeURL = BACKEND_URL + '/' + removeURL;
			}
			if (postField == null) postField = 'id';
			var data = [];
			$(elem).each(function () {
				if ($(this).is(':checked') && !$(this).prop('disabled')) data.push({ name: postField, value: $(this).val() });
			});
			if (data.length < 1) {
				App.message({ title: App.i18n.SYS_INFO, text: unselectedMsg || App.i18n.PLEASE_SELECT_FOR_OPERATE, type: 'warning' });
				return false;
			}

			var answer = confirmMsg;
			if (answer) {
				answer += "\n" + App.sprintf(App.i18n.SELECTED_ITEMS, data.length);
				if (!confirm(answer)) return false;
			}
			App.loading('show');
			$.get(removeURL, data, function (r) {
				App.loading('hide');
				if (callback && typeof callback === "function") return callback();
				var msg = { title: App.i18n.SYS_INFO, text: r.Info, type: '' };
				if (r.Code == 1) {
					msg.type = 'success';
					if (!msg.text) msg.text = '操作成功';
				} else {
					msg.type = 'error';
					if (!msg.text) msg.text = '操作失败';
				}
				App.message(msg);
				window.setTimeout(function () {
					window.location.reload();
				}, 2000);
			}, 'json');
			return true;
		},
		parseBool: function (b) {
			switch (String(b).toLowerCase()) {
				case '0':
				case 'false':
				case 'n':
				case 'no':
				case 'off':
				case 'null':
					return false;

				default:
					return true;
			}
		},
		progressMonitor: function (getCurrentFn, totalProgress) {
			NProgress.start();
			var interval = window.setInterval(function () {
				var current = getCurrentFn() / totalProgress;
				if (current >= 1) {
					NProgress.set(1);
					window.clearInterval(interval);
				} else {
					NProgress.set(current);
				}
			}, 50);
		},
		floatAtoi: function(p, defaults){
			var mp = {
				'bottom':'7',
				'right':'6',
				'top':'5',
				'left':'8',
				'leftBottom':'4',
				'rightBottom':'3',
				'leftTop':'1',
				'rightTop':'2',
			};
			if(typeof mp[p] !== 'undefined') return mp[p];
			return defaults;
		},
		float: function (elem, mode, attr, position, options) {
			if (!mode) mode = 'ajax';
			if (!attr) attr = mode=='remind'?'rel':'src';
			if (!position) position = '5-7';//两个数字分别代表trigger(触发对象)-target(浮动层)，（各个数字的编号从矩形框的左上角开始，沿着顺时针开始旋转来进行编号，然后再从上中部开始沿着顺时针开始编号进行。也就是1、2、3、4分别代表左上角、右上角、右下角、左下角；5、6、7、8分别代表上中、右中、下中、左中）
			else {
				var arr = String(position).split('-');
				var val = [];
				for (var i = 0; i < arr.length && i < 2; i++) {
					var p = arr[i];
					val.push(App.floatAtoi(p, p));
				}
				if (arr.length==1) {
					switch (arr[0]) {
						case 'bottom':val.push(App.floatAtoi('top'));break;
						case 'right':val.push(App.floatAtoi('left'));break;
						case 'top':val.push(App.floatAtoi('bottom'));break;
						case 'left':val.push(App.floatAtoi('right'));break;
						case 'leftBottom':val.push(App.floatAtoi('rightBottom'));break;
						case 'rightBottom':val.push(App.floatAtoi('leftBottom'));break;
						case 'leftTop':val.push(App.floatAtoi('rightTop'));break;
						case 'rightTop':val.push(App.floatAtoi('leftTop'));break;
						default:val.push(p);break;
					}
				}
				position = val.join('-');
			}
			//console.log(position);
			var defaults = { 'targetMode': mode, 'targetAttr': attr, 'position': position };
			$(elem).powerFloat($.extend(defaults,options||{}));
		},
		uploadPreviewer: function (elem, options, successCallback, errorCallback) {
			if($(elem).parent('.file-preview-shadow').length<1){
				var defaults = {
					"buttonText":'<i class="fa fa-cloud-upload"></i> '+App.i18n.BUTTON_UPLOAD,
					"previewTableContainer":'#previewTableContainer',
					"url":'',
					"previewTableShow":false,
					"uploadProgress":function(progress){
						var count=progress*100;
						if(count>100){
							$.LoadingOverlay("hide");
							return;
						}
						$.LoadingOverlay("progress", count);
					}
				};
				var noptions = $.extend({}, defaults, options || {});
				var uploadInput = $(elem).uploadPreviewer(noptions);
				$(elem).data('uploadPreviewer', uploadInput);
				$(elem).on("file-preview:changed", function(e) {
					var options = {
						image: ASSETS_URL+"/images/nging-gear.png",
						progress: false, 
						//fontawesome : "fa fa-cog fa-spin",
						text: App.i18n.UPLOADING
					};
					if(noptions.uploadProgress){
						options.progress = true;
						options.image = "";
					}
				  	$.LoadingOverlay("show", options);
				  	uploadInput.submit(function(r){
					  $.LoadingOverlay("hide");
					  if(r.Code==1){
						  App.message({text:App.i18n.UPLOAD_SUCCEED,type:'success'});
						  if(successCallback!=null) successCallback(r);
					  }else{
						  App.message({text:r.Info,type:'error'});
						  if(errorCallback!=null) errorCallback(r);
					  }
				  	},function(){
						if(errorCallback!=null) errorCallback();
					});
				});
			}
		},
		showRequriedInputStar:function(){
			$('form:not([required-redstar])').each(function(){
				$(this).attr('required-redstar','1');
				$(this).find('[required]').each(function(){
					var parent = $(this).parent('.input-group');
					if(parent.length>0){
						if(!parent.hasClass('required')) parent.addClass('required');
						return;
					}
					parent = $(this).closest('.form-group,div[class*="col-"]');
					if(parent.length<1) return;
					var lbl;
					if(parent.hasClass('form-group')){
						lbl = parent.children('.control-label');
					}else{
						lbl = parent.prev('.control-label');
					}
					if (lbl.length<1 || lbl.hasClass('required')) return;
					lbl.addClass('required');
				});
			});
		},
		pushState:function(data,title,url){
			if(!window.history || !window.history.pushState)return;
			window.history.pushState(data,title,url);
		},
		replaceState:function(data,title,url){
			if(!window.history || !window.history.replaceState)return;
			window.history.replaceState(data,title,url);
		},
		formatJSON:function(json){
			json = $.trim(json);
			var first = json.substring(0,1);
			if (first=='['||first=='{'){
				var obj = JSON.parse(json);
				json = JSON.stringify(obj, null, "\t");
				return json;
			}
			return '';
		},
		formatJSONFromInnerHTML:function($a){
			if($a.data('jsonformatted'))return;
			$a.data('jsonformatted',true);
			var json = App.formatJSON($a.html());
			if(json!='') $a.html(json);
		},
		rewriteFormValueAsArray: function(formData,fieldName,arrayData){
			var j = 0;
			for (var i=0; i<formData.length; i++) {
				var v = formData[i];
				if(v.name == fieldName) {
					if(arrayData.length > j){
						v.value = arrayData[j];
						j++;
					}else{
						v.value = '';
					}
				}
			}
			for (var i=j; i<arrayData.length; i++) {
				formData.push({name:fieldName,value:arrayData[i]});
			}
			return formData;
		},
		plotTooltip: function (x, y, contents) {
			$("<div id='plot-tooltip'>" + contents + "</div>").css({
			  position: "absolute",
			  display: "none",
			  top: y + 5,
			  left: x + 5,
			  border: "1px solid #000",
			  padding: "5px",
			  'color':'#fff',
			  'border-radius':'2px',
			  'font-size':'11px',
			  "background-color": "#000",
			  opacity: 0.80
			}).appendTo("body").fadeIn(200);
		},
		plotHover: function(elem, formatter) {
			$(elem).on("plothover", function (event, pos, item) { //var str = "(" + pos.x.toFixed(2) + ", " + pos.y.toFixed(2) + ")";
			  	if (!item) {
			    	$("#plot-tooltip").remove();
			    	previousPlotPoint = null;
					return;
			  	}
				if (previousPlotPoint == item.dataIndex) return;
				previousPlotPoint = item.dataIndex;
			    $("#plot-tooltip").remove();
				if(formatter == null) formatter = function(event, pos, item) {
					var x = item.datapoint[0].toFixed(2), y = item.datapoint[1].toFixed(2);
					return item.series.label + " of " + x + " = " + y;
				}
			    App.plotTooltip(item.pageX, item.pageY, formatter.call(this, event, pos, item));
			}); 
		},
		template: function(elem,jsonData,onSwitchPage){
            var tplId = $(elem).attr('tpl'), arr = String(tplId).split('=>');
			tplId = $.trim(arr[0]).replace(/^#/,'');
			if($('#'+tplId).length<1) throw new Error('not found template: '+(typeof elem == 'string' ? elem : $(elem).attr('id'))+'=>'+tplId);
			if(arr.length > 1){
				elem = $.trim(arr[1]);
				if($(elem).length<1) throw new Error('[template] not found wrapper: '+elem);
			}
            $(elem).html(template(tplId,jsonData));
			if(onSwitchPage==null) return;
			$(elem).find('ul.pagination li > a[page]').on('click',function(){
				if($(this).closest('li.disabled').length > 0) return;
				onSwitchPage($(this).attr('page'));
			})
			$(elem).find('ul.pagination').on('refresh',function(){
				var page=1;
				if($(this).data('page')){
					page=$(this).data('page');
					$(this).data('page',false);
				}else{
					page=$(this).find('li.active > a[page]').data('page')||1;
				}
				onSwitchPage(page);
			});
		},
		pagination: function(data, num){
			var pageNumbers = [], page = data.page, pages = data.pages;
			if(!num) num = 10;
			var remainPages = pages - page, start = 1, count = 0;
			if(remainPages < num){
				start = pages - num + 1
			} else {
				start = page - (num / 2)
			}
			if(start < 1) start = 1;
			for(var i = start; i <= pages; i++){
				count++
				if(count > num) break;
				pageNumbers.push(i);
			}
			var pagination = {
				size:data.size,
				page:page,
				rows:data.rows,
				pages:pages,
				pageNumbers:pageNumbers,
			}
			return pagination;
		},
		withPagination: function(data, num, tpl){
			var pagination = App.pagination(data, num);
			if(tpl==null) tpl = 'tpl-pagination';
			data.pagination=template(tpl,pagination);
			return data;
		},
		captchaUpdate: function($form, resp){
			if(App.captchaHasError(resp.Code) && resp.Data && typeof(resp.Data.captchaIdent) !== 'undefined') {
				if(false == ($form instanceof jQuery)) $form=$($form);
				var idElem = $form.find('input#'+resp.Data.captchaIdent);
				idElem.val(resp.Data.captchaID);
				idElem.siblings('img').attr('src',resp.Data.captchaURL);
				if(resp.Data.captchaName) $form.find('input[name="'+resp.Data.captchaName+'"]').focus();
			}
		},
		captchaHasError: function(code) {
			return code >= -11 && code <= -9;
		},
		treeToggle: function (elem, options) {
			var defaults = {
				expandFirst: false,
				ajax: null,
				onclick: null,
			};
			options = $.extend(defaults, options||{});
			if(!elem) elem = '.treeview';
			var expand = function(){
				var icon = $(this).children(".fa"), tree;
				if(icon.length > 0){
					tree = $(this).next('ul.tree');
				} else if($(this).hasClass('fa')){
					icon = $(this);
					tree = $(this).parent().next('ul.tree');
				} else {
					return;
				}
				if(icon.hasClass("fa-folder-o")){
				  icon.removeClass("fa-folder-o").addClass("fa-folder-open-o");
				}else{
				  icon.removeClass("fa-folder-open-o").addClass("fa-folder-o");
				}
				tree.toggle(300,function(){
				  $(this).parent().toggleClass("open");
				  $(".tree .nscroller").nanoScroller({ preventPageScrolling: true });
				  $(window).trigger('resize');
				});
			};
			$(elem).on('click','.tree-toggler',function () {
				var that = this;
				if(typeof options.onclick === "function") options.onclick.apply(that, arguments);
				if($(that).data('loaded')||!options.ajax) return expand.apply(that);
				$(that).data('loaded',true);
				var ajaxOptions = options.ajax;
				if(typeof options.ajax === "function"){
					ajaxOptions = options.ajax.apply(that, arguments);
				}
				if(typeof ajaxOptions.data === "function"){
					ajaxOptions.data = ajaxOptions.data.apply(that, arguments);
				}
				$.ajax(ajaxOptions).done(function(){
					expand.apply(that,arguments);
				});
			})
			if(options.expandFirst) $(elem+' label.tree-toggler:first').trigger('click');
		},
		isFunction: function(v){
			return typeof v === "function";
		},
		currentURL: function(){
  			if(typeof(window.IS_BACKEND)!='undefined' && window.IS_BACKEND){
    			return BACKEND_URL;
  			}
  			return FRONTEND_URL;
		}
	};

}();

$(function () {
	//$("body").animate({opacity:1,'margin-left':0},500);
	$("body").css({ opacity: 1, 'margin-left': 0 });
});
