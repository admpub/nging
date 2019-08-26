/**
* datetimepicker快捷函数
* @author Hank Shen <swh@admpub.com>
* @param {string|object} startTimeElem 
* @param {string|object|null} endTimeElem 
* @param {object} options
*/
App.datetimepicker=function(startTimeElem,endTimeElem,options){
 var config = {
   'language': 'zh-CN', 
   'format': 'yyyy-mm-dd',
   'minView': 'month',
   'todayBtn': 1, 
   'autoclose': 1,
   'todayHighlight': 1
 };
 config = $.extend(config,options||{});
 var startTime=$(startTimeElem).datetimepicker(config);
 if(!endTimeElem) return;
 var endTime=$(endTimeElem).datetimepicker(config);
 startTime.on('changeDate', function (e) {
   endTime.datetimepicker("setStartDate", e.date);
 });
 endTime.on('changeDate', function (e) {
   startTime.datetimepicker("setEndDate", e.date);
 });
};
