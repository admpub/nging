/**
* datetimepicker快捷函数
* @author Hank Shen <swh@admpub.com>
* @param {string|object} startTimeElem 
* @param {string|object|null} endTimeElem 
* @param {object} options
*/
App.datetimepicker = function (startTimeElem, endTimeElem, options) {
	var config = {
		'language': 'zh-CN',
		'format': 'yyyy-mm-dd',
		'minView': 'month',
		'todayBtn': true,
		'autoclose': true,
		'todayHighlight': true
	};
	config = $.extend(config, options || {});
	var startTime = $(startTimeElem).datetimepicker(config);
	if (!endTimeElem) return startTime;
	var endTime = $(endTimeElem).datetimepicker(config);
	startTime.on('changeDate', function (e) {
		endTime.datetimepicker("setStartDate", e.date);
	});
	endTime.on('changeDate', function (e) {
		startTime.datetimepicker("setEndDate", e.date);
	});
	if ($(startTimeElem).val()) endTime.datetimepicker('setStartDate', $(startTimeElem).val());
	if ($(endTimeElem).val()) startTime.datetimepicker('setEndDate', $(endTimeElem).val());
	return [startTime,endTime];
};
/**
 * 日期范围选择
 * @denpend js/daterangepicker/daterangepicker.min.css
 * @denpend js/daterangepicker/moment.min.js
 * @denpend js/daterangepicker/jquery.daterangepicker.min.js
 * @document https://longbill.github.io/jquery-date-range-picker/
 */
App.daterangepicker = function (rangeElem, options) {
	var change = false;
	var config = {
		customArrowPrevSymbol: '<i class="fa fa-arrow-left"></i>',
		customArrowNextSymbol: '<i class="fa fa-arrow-right"></i>',
		autoClose: true,
		format: 'YYYY-MM-DD',//YYYY-MM-DD HH:mm:ss
		separator: ' - ',
		singleDate: false,
		language:'cn',
		monthSelect: true,
		maxDays: 30,
		yearSelect: true, //[1900, moment().get('year')]
		time: {
			enabled: false
		}
	};
	config = $.extend(config, options || {});
	if(!config.time.enabled && config.format.indexOf('H')!=-1) config.time.enabled = true;
	$(rangeElem).on('focus click touch',function(){
		$(this).select();
	});
	return $(rangeElem).dateRangePicker(config).bind('datepicker-closed',function(){
		$(rangeElem).focus();
		if(change) $(rangeElem).trigger('change');
	}).bind('datepicker-change',function(){
		change = true;
	});
	//$(rangeElem).data('dateRangePicker').setDateRange('2013-11-20','2013-11-25');
};
App.daterangepickerx = function(container,startElement,endElement,options){
	var defaults = {
		getValue: function(){
		  if ($(startElement).val() && $(endElement).val() ) return $(startElement).val() + ' - ' + $(endElement).val();
		  return '';
		},
		setValue: function(s,s1,s2){
		  $(startElement).val(s1);
		  $(endElement).val(s2);
		},
	};
	App.daterangepicker(container,$.extend(defaults,options||{}));
};
App.datepicker = function (elem, options) {
	if (!options) options = {};
	options.singleDate = true;
	options.singleMonth = true;
	options.showShortcuts = false;
	options.showTopbar = false;
	return App.daterangepicker(elem, options);
};