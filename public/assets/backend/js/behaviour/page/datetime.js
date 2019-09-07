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
	if (!endTimeElem) return;
	var endTime = $(endTimeElem).datetimepicker(config);
	startTime.on('changeDate', function (e) {
		endTime.datetimepicker("setStartDate", e.date);
	});
	endTime.on('changeDate', function (e) {
		startTime.datetimepicker("setEndDate", e.date);
	});
	if ($(startTimeElem).val()) endTime.datetimepicker('setStartDate', $(startTimeElem).val());
	if ($(endTimeElem).val()) startTime.datetimepicker('setEndDate', $(endTimeElem).val());
};
/**
 * 日期范围选择
 * @denpend js/daterangepicker/daterangepicker.min.css
 * @denpend js/daterangepicker/moment.min.js
 * @denpend js/daterangepicker/jquery.daterangepicker.min.js
 */
App.daterangepicker = function (rangeElem, options) {
	var config = {
		autoClose: true,
		format: 'YYYY-MM-DD',
		separator: ' - ',
		singleDate: false,
		language:'cn',
		monthSelect: true,
    	yearSelect: true //[1900, moment().get('year')]
	};
	config = $.extend(config, options || {});
	$(rangeElem).dateRangePicker(config);
	//$(rangeElem).data('dateRangePicker').setDateRange('2013-11-20','2013-11-25');
};
App.datepicker = function (elem, options) {
	if (!options) options = {};
	options.singleDate = true;
	options.singleMonth = true;
	options.showShortcuts = false;
	options.showTopbar = false;
	App.daterangepicker(elem, options);
};